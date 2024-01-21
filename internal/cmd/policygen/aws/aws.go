package aws

import (
	"fmt"

	parser "github.com/nukleros/markers"
	"github.com/scottd018/go-utils/pkg/directory"
	"github.com/spf13/cobra"

	"github.com/scottd018/policy-gen/internal/pkg/aws"
	"github.com/scottd018/policy-gen/internal/pkg/input"
)

const awsPolicyGenExample = `
policygen aws --input-path=<input_path> --output-path='<output_path>'
`

const (
	flagInputPath  = "input-path"
	flagOutputPath = "output-path"

	markerDefinition = "+aws:iam:policy"
)

type awsPolicyGenInputs struct {
	inputPath  string
	outputPath string
}

func NewCommand() *cobra.Command {
	// add a place to store user input from the command
	input := awsPolicyGenInputs{}

	// create the command
	command := &cobra.Command{
		Use:     "aws",
		Short:   "Generate AWS IAM policies",
		Long:    `Generate AWS IAM policies`,
		PreRunE: func(cmd *cobra.Command, args []string) error { return setup(input) },
		RunE:    func(cmd *cobra.Command, args []string) error { return run(input) },
		Example: awsPolicyGenExample,
	}

	// add flags
	command.Flags().StringVarP(&input.inputPath, flagInputPath, "i", "./", "Input path to recursively begin parsing markers")
	command.Flags().StringVarP(&input.outputPath, flagOutputPath, "o", "./", "Output path to output generated policies")

	return command
}

func setup(inputs awsPolicyGenInputs) error {
	directoryInputs := map[string]string{
		flagInputPath:  inputs.inputPath,
		flagOutputPath: inputs.outputPath,
	}

	// ensure our inputs are not empty
	for flag, value := range directoryInputs {
		if value == "" {
			return fmt.Errorf("missing value for required flag [%s]", flag)
		}
	}

	// ensure our directory inputs point at an existing path
	for flag, path := range directoryInputs {
		exists, err := directory.Exists(path)
		if err != nil {
			return fmt.Errorf("directory path [%s] for flag [%s] is invalid - %w", path, flag, err)
		}

		if !exists {
			return fmt.Errorf("directory path [%s] for flag [%s] is missing", path, flag)
		}
	}

	return nil
}

func run(inputs awsPolicyGenInputs) error {
	policyMarker := aws.Marker{}

	// create a registry for our field markers
	registry := parser.NewRegistry()

	// define our marker
	definition, err := parser.Define(markerDefinition, policyMarker)
	if err != nil {
		return fmt.Errorf("unable to create policy definition for marker [%s] - %w", markerDefinition, err)
	}

	// add the marker to the registry
	registry.Add(definition)

	// collect the data from the given path
	data, err := input.Collect(inputs.inputPath)
	if err != nil {
		return fmt.Errorf("unable to collect file data to be parsed from path [%s] - %w", inputs.inputPath, err)
	}

	// run the parser
	results := parser.NewParser(string(data), registry).Parse()
	if len(results) == 0 {
		fmt.Printf("no results found for marker [%s] at path [%s]\n", markerDefinition, inputs.inputPath)

		return nil
	}

	// process the results
	_, err = aws.FindMarkers(results)
	if err != nil {
		return fmt.Errorf("unable to convert results to markers - %w", err)
	}

	return nil
}
