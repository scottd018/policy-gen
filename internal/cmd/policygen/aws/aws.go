package aws

import (
	"fmt"

	"github.com/scottd018/go-utils/pkg/directory"
	"github.com/spf13/cobra"

	"github.com/scottd018/policy-gen/internal/pkg/aws"
	"github.com/scottd018/policy-gen/internal/pkg/input"
)

const awsPolicyGenExample = `
policygen aws --input-path=<input_path> --output-path='<output_path>' --force
`

func NewCommand() *cobra.Command {
	// add a place to store user in from the command
	in := input.Input{}

	// create the command
	command := &cobra.Command{
		Use:     "aws",
		Short:   "Generate AWS IAM policies",
		Long:    `Generate AWS IAM policies`,
		PreRunE: func(cmd *cobra.Command, args []string) error { return setup(in) },
		RunE:    func(cmd *cobra.Command, args []string) error { return run(in) },
		Example: awsPolicyGenExample,
	}

	// add flags
	command.Flags().StringVarP(
		&in.InputPath, input.FlagInputPath, input.FlagInputPathShort, input.FlagInputPathDefault,
		input.FlagInputPathDescription,
	)

	command.Flags().StringVarP(
		&in.OutputPath, input.FlagOutputPath, input.FlagOutputPathShort, input.FlagOutputPathDefault,
		input.FlagOutputPathDescription,
	)

	command.Flags().BoolVarP(
		&in.Force, input.FlagForce, input.FlagForceShort, input.FlagForceDefault,
		input.FlagForceDescription,
	)

	command.Flags().BoolVar(
		&in.Debug, input.FlagDebug, input.FlagDebugDefault,
		input.FlagDebugDescription,
	)

	return command
}

func setup(in input.Input) error {
	directoryInputs := map[string]string{
		input.FlagInputPath:  in.InputPath,
		input.FlagOutputPath: in.OutputPath,
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

func run(inputs input.Input) error {
	processor := aws.NewMarkerProcessor(inputs)
	if err := processor.Process(); err != nil {
		return fmt.Errorf("unable to process markers - %w", err)
	}

	return nil
}
