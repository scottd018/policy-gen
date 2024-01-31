package aws

import (
	"fmt"
	"os"

	"github.com/nukleros/markers"
	"github.com/nukleros/markers/parser"
	"github.com/rs/zerolog"

	"github.com/scottd018/policy-gen/internal/pkg/input"
)

// MarkerProcessor represents the object used to process markers
// for a file.
type MarkerProcessor struct {
	Input input.Input
	Log   zerolog.Logger
}

// NewMarkerProcessor instantiates a new instance of a markerProcessor
// object.
func NewMarkerProcessor(inputs input.Input) *MarkerProcessor {
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
		return fmt.Errorf("error parsing marker results from input path [%s] - %w", processor.Input.InputPath, err)
	}

	// convert our file markers into aws markers
	awsMarkers, err := processor.FindMarkers(results)
	if err != nil {
		return fmt.Errorf("error converting results to markers - %w", err)
	}

	// process the markers and write the files
	for file, document := range awsMarkers.PolicyFiles() {
		filename := Filename(processor.Input.OutputPath, file)

		processor.Log.Info().Msgf("writing file: [%s]", filename)

		if err := document.Write(filename, processor.Input.Force); err != nil {
			return err
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
	definition, err := markers.Define(MarkerDefinition(), policyMarker)
	if err != nil {
		processor.Log.Info().Msgf("adding marker definition to registry: [%s]", MarkerDefinition())

		return nil, fmt.Errorf("unable to create policy definition for marker [%s] - %w", MarkerDefinition(), err)
	}

	// add the marker to the registry
	registry.Add(definition)

	// collect the data from the given path
	data, err := input.Collect(processor.Input.InputPath)
	if err != nil {
		processor.Log.Info().Msgf("collecting input for path: [%s]", processor.Input.InputPath)

		return nil, fmt.Errorf("error collecting file data for marker [%s] - %w", MarkerDefinition(), err)
	}

	// run the parser
	results := markers.NewParser(string(data), registry).Parse()
	if len(results) == 0 {
		processor.Log.Warn().Msgf("no results found for marker [%s] at path [%s]\n", MarkerDefinition(), processor.Input.InputPath)

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
