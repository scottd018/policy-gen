package aws

import (
	"fmt"

	"github.com/scottd018/go-utils/pkg/directory"
	"github.com/spf13/cobra"

	"github.com/scottd018/policy-gen/internal/pkg/aws"
)

const awsPolicyGenExample = `
policygen aws --input-path=<input_path> --output-path='<output_path>' --force
`

const (
	flagInputPath  = "input-path"
	flagOutputPath = "output-path"
	flagForce      = "force"
)

type awsPolicyGenInputs struct {
	inputPath  string
	outputPath string
	force      bool
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
	command.Flags().BoolVarP(&input.force, flagForce, "f", false, "Forcefully overwrite files with matching names")

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
	// retrieve the marker results from the input path
	results, err := aws.MarkerResults(inputs.inputPath)
	if err != nil {
		return fmt.Errorf("error finding marker results from input path [%s] - %w", inputs.inputPath, err)
	}

	// convert our file markers into aws markers
	awsMarkers, err := aws.FindMarkers(results)
	if err != nil {
		return fmt.Errorf("error converting results to markers - %w", err)
	}

	// process the markers and write the files
	err = awsMarkers.Process().Write(inputs.outputPath, inputs.force)
	if err != nil {
		return fmt.Errorf("error writing files to output path [%s] - %w", inputs.outputPath, err)
	}

	return nil
}
