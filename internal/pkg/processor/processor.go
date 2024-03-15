package processor

import (
	"fmt"
	"os"
	"unicode/utf8"

	"github.com/nukleros/markers"
	"github.com/nukleros/markers/marker"
	"github.com/nukleros/markers/parser"
	"github.com/rs/zerolog"

	"github.com/scottd018/policy-gen/internal/pkg/docs"
	"github.com/scottd018/policy-gen/internal/pkg/files"
	"github.com/scottd018/policy-gen/internal/pkg/policy"
	"github.com/scottd018/policy-gen/internal/pkg/utils"
)

// Processor represents the object used to process markers
// for a file.
type Processor struct {
	Config              *Config
	Log                 zerolog.Logger
	Definition          *marker.Definition
	Registry            *marker.Registry
	PolicyFileGenerator policy.FileGenerator
}

// NewProcessor instantiates a new instance of a Processor object.  A processor
// is used to process a given set of markers from a given set of inputs, mainly
// the input path to parse.
func NewProcessor(config *Config, marker string, object interface{}, generator policy.FileGenerator) (*Processor, error) {
	// configure logging
	level := zerolog.InfoLevel
	if config.Debug {
		level = zerolog.DebugLevel
	}

	logger := zerolog.ConsoleWriter{
		Out: os.Stdout,
		PartsExclude: []string{
			"time",
		},
	}

	// create a registry for our field markers
	registry := markers.NewRegistry()

	// define our marker
	definition, err := markers.Define(marker, object)
	if err != nil {
		return nil, fmt.Errorf("unable to create policy definition for marker [%s] - %w", marker, err)
	}

	// add the marker to the registry
	registry.Add(definition)

	return &Processor{
		Config:              config,
		Log:                 zerolog.New(logger).With().Timestamp().Logger().Level(level),
		Registry:            registry,
		Definition:          definition,
		PolicyFileGenerator: generator,
	}, nil
}

// Process executes the marker processing.
func (processor *Processor) Process() error {
	// retrieve the marker results from the input path
	results, err := processor.Parse()
	if err != nil {
		return fmt.Errorf(
			"error parsing marker results from input path [%s] - %w",
			processor.Config.InputDirectory.Path,
			err,
		)
	}

	// convert our file markers into a set of policyMarkers markers
	policyMarkers, err := processor.FindMarkers(results)
	if err != nil {
		return fmt.Errorf("error converting results to markers - %w", err)
	}

	// retrieve our policy files from our markers
	policyFiles, err := policy.ToPolicyFiles(policyMarkers, processor.PolicyFileGenerator)
	if err != nil {
		return fmt.Errorf("error retrieving files from markers - %w", err)
	}

	// write our files
	options := []files.Option{}
	if processor.Config.Force {
		options = []files.Option{files.WithOverwrite}
	}

	for _, policyFile := range policyFiles {
		processor.Log.Info().Msgf("writing policy file: [%s]", policyFile.Path())
		if err := policyFile.Write(files.ModePolicyFile, options...); err != nil {
			return fmt.Errorf("error writing policy file: [%s] - %w", policyFile.Path(), err)
		}
	}

	// write the documentation if it was requested
	if processor.Config.DocumentationFile != nil && processor.Config.DocumentationFile.File != "" {
		processor.Log.Info().Msgf("writing documentation file: [%s]", processor.Config.DocumentationFile.Path())

		// create the document and generate the content
		documentationFile := docs.NewDocumentation(processor.Config.DocumentationFile)
		documentationFile.Generate(ToDocumentRows(policyMarkers)...)

		// write the documentation to the specified path
		if err := documentationFile.File.Write(files.ModePolicyFile, options...); err != nil {
			return fmt.Errorf("error writing documentation file: [%s] - %w", documentationFile.File.Path(), err)
		}
	}

	return nil
}

// Parse parses a set of markers from a given path and returns the results.
func (processor *Processor) Parse() ([]*parser.Result, error) {
	processor.Log.Info().Msgf("parsing markers: [%s]", processor.Definition.Name)
	processor.Log.Info().Msgf("collecting input for path: [%s]", processor.Config.InputDirectory.Path)

	// collect the file paths from the given input directory path
	files, err := processor.Config.InputDirectory.ListFilePaths(processor.Config.Recursive)
	if err != nil {
		return nil, fmt.Errorf("error collecting file paths for marker: [%s] - %w", processor.Definition.Name, err)
	}

	// parse the content of each file and collect the results
	results := []*parser.Result{}

	for path := range files {
		fullPath := fmt.Sprintf("%s/%s", processor.Config.InputDirectory.Path, files[path])

		processor.Log.Debug().Msgf("collecting marker results for file: [%s]", fullPath)

		// read in the file content
		content, err := os.ReadFile(fullPath)
		if err != nil {
			return nil, fmt.Errorf("unable to read file: [%s] - %w", fullPath, err)
		}

		// only append text file content
		if utf8.Valid(content) {
			found := markers.NewParser(string(content), processor.Registry).Parse()
			if len(found) != 0 {
				results = append(results, found...)
			}
		}
	}

	// warn the user if we found no markers based on their input
	if len(results) == 0 {
		processor.Log.Warn().Msgf(
			"no results found for marker [%s] at path: [%s]\n",
			processor.Definition.Name,
			processor.Config.InputDirectory.Path,
		)

		return []*parser.Result{}, nil
	}

	return results, nil
}

// FindMarkers finds all the markers in a given set of parsed results.
func (processor *Processor) FindMarkers(results []*parser.Result) ([]policy.Marker, error) {
	foundMarkers := make([]policy.Marker, len(results))

	for i := range results {
		// convert the marker to its underlying type
		marker, err := utils.ConvertToMarker(results[i].Object)
		if err != nil {
			return nil, fmt.Errorf(
				"found invalid marker with text [%s] at position [%d] - %w",
				results[i].MarkerText,
				i,
				err,
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

// ToDocumentRows converts a Markers object to a set of document row interfaces.  This is needed
// to display markers in documentation.
func ToDocumentRows(m []policy.Marker) []docs.Row {
	markersSlice := make([]docs.Row, len(m))

	for i := range m {
		markersSlice[i] = m[i]
	}

	return markersSlice
}
