package aws

import (
	"fmt"
	"os"

	"github.com/nukleros/markers"
	"github.com/nukleros/markers/parser"
	"github.com/rs/zerolog"

	"github.com/scottd018/policy-gen/internal/pkg/docs"
	"github.com/scottd018/policy-gen/internal/pkg/files"
	"github.com/scottd018/policy-gen/internal/pkg/input"
)

// MarkerProcessor represents the object used to process markers
// for a file.
type MarkerProcessor struct {
	Input *input.Processor
	Log   zerolog.Logger
}

// NewMarkerProcessor instantiates a new instance of a markerProcessor
// object.
func NewMarkerProcessor(inputs *input.Processor) *MarkerProcessor {
	level := zerolog.InfoLevel
	if inputs.Debug {
		level = zerolog.DebugLevel
	}

	return &MarkerProcessor{
		Input: inputs,
		Log:   zerolog.New(os.Stdout).With().Timestamp().Logger().Level(level),
	}
}

// Process executes the marker processing.
func (processor *MarkerProcessor) Process() error {
	// retrieve the marker results from the input path
	results, err := processor.Parse()
	if err != nil {
		return fmt.Errorf("error parsing marker results from input path [%s] - %w", processor.Input.InputDirectory.Path, err)
	}

	// convert our file markers into aws markers
	awsMarkers, err := processor.FindMarkers(results)
	if err != nil {
		return fmt.Errorf("error converting results to markers - %w", err)
	}

	// process the markers and write the files
	for policyFile, document := range awsMarkers.PolicyFiles() {
		rawFileName := PolicyFilename(processor.Input.OutputDirectory.Path, policyFile)

		// we do not need to pass the pre-existing directory option here because it
		// was validated on input
		jsonFile, err := files.NewJSONFile(rawFileName)
		if err != nil {
			return fmt.Errorf("error creating policy file object [%s] - %w", policyFile, err)
		}

		processor.Log.Info().Msgf("writing policy file: [%s]", jsonFile.Path())

		if err := document.Write(jsonFile, processor.Input.Force); err != nil {
			return fmt.Errorf("error writing policy file [%s] - %w", jsonFile.Path(), err)
		}
	}

	// write the documentation if it was requested
	if processor.Input.DocumentationFile != nil && processor.Input.DocumentationFile.File != "" {
		// create the documentation file
		documentationFile, err := files.NewMarkdownFile(processor.Input.DocumentationFile.File)
		if err != nil {
			return fmt.Errorf("error creating markdown file object [%s] - %w", processor.Input.DocumentationFile.File, err)
		}

		processor.Log.Info().Msgf("writing documentation file: [%s]", documentationFile.Path())

		// create the document
		if err := docs.NewDocumentation(documentationFile).Write(
			processor.Input.Force,
			awsMarkers.ToDocumentRows()...,
		); err != nil {
			return fmt.Errorf("error writing documentation file [%s] - %w", documentationFile.Path(), err)
		}
	}

	return nil
}

// Parse parses a set of markers from a given path and returns the results.
func (processor *MarkerProcessor) Parse() ([]*parser.Result, error) {
	policyMarker := Marker{}

	// create a registry for our field markers
	registry := markers.NewRegistry()

	// define our marker
	processor.Log.Info().Msgf("parsing markers: [%s]", MarkerDefinition())
	definition, err := markers.Define(MarkerDefinition(), policyMarker)
	if err != nil {
		return nil, fmt.Errorf("unable to create policy definition for marker [%s] - %w", MarkerDefinition(), err)
	}

	// add the marker to the registry
	registry.Add(definition)

	// collect the data from the given path
	processor.Log.Info().Msgf("collecting input for path: [%s]", processor.Input.InputDirectory.Path)
	data, err := processor.Input.InputDirectory.CollectData()
	if err != nil {
		return nil, fmt.Errorf("error collecting file data for marker [%s] - %w", MarkerDefinition(), err)
	}

	// run the parser
	results := markers.NewParser(string(data), registry).Parse()
	if len(results) == 0 {
		processor.Log.Warn().Msgf("no results found for marker [%s] at path [%s]\n", MarkerDefinition(), processor.Input.InputDirectory.Path)

		return []*parser.Result{}, nil
	}

	return results, nil
}

// FindMarkers finds all the markers in a given set of parsed results.
func (processor *MarkerProcessor) FindMarkers(results []*parser.Result) (Markers, error) {
	foundMarkers := make(Markers, len(results))

	for i := range results {
		// ensure the marker we found is the appropriate object type
		marker, ok := results[i].Object.(Marker)
		if !ok {
			return nil, fmt.Errorf(
				"found invalid marker with text [%s] at position [%d]",
				results[i].MarkerText,
				i,
			)
		}

		// ensure the marker we found is valid
		if err := marker.Validate(); err != nil {
			return nil, fmt.Errorf(
				"found invalid marker with text [%s] at position [%d] - %w",
				results[i].MarkerText,
				i,
				err,
			)
		}

		processor.Log.Debug().Msgf("found marker: [%s]", results[i].MarkerText)

		// add the markers to the slice
		foundMarkers[i] = marker
	}

	return foundMarkers, nil
}
