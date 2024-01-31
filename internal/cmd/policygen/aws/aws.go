package aws

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/scottd018/policy-gen/internal/pkg/aws"
	"github.com/scottd018/policy-gen/internal/pkg/input"
)

const awsPolicyGenExample = `
# generate policies using sensible defaults
policygen aws

# generate policies from files at located at input path ./input and write 
# discovered policies to output path ./output while forcefully overwriting
# any overlapping policies in the ./output directory.
policygen aws --input-path=./input --output-path=./output --force

# generate policies with debug logging
policygen aws --debug

# generate policies and associated documentation at ./output/README.md
policygen aws --output-path=./output --documentation=README.md
`

func NewCommand() *cobra.Command {
	flags := input.NewFlags()

	// create the command
	command := &cobra.Command{
		Use:     "aws",
		Short:   "Generate AWS IAM policies",
		Long:    `Generate AWS IAM policies`,
		PreRunE: func(cmd *cobra.Command, args []string) error { return setup(flags) },
		RunE:    func(cmd *cobra.Command, args []string) error { return run(flags) },
		Example: awsPolicyGenExample,
	}

	// add flags
	flags.Initialize(command)

	return command
}

func setup(flags input.Flags) error {
	return nil
}

func run(flags input.Flags) error {
	processorInputs, err := flags.Process()
	if err != nil {
		return fmt.Errorf("unable to process flags - %w", err)
	}

	processor := aws.NewMarkerProcessor(processorInputs)
	if err := processor.Process(); err != nil {
		return fmt.Errorf("unable to process markers - %w", err)
	}

	return nil
}
