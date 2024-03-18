package aws

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/scottd018/policy-gen/internal/pkg/aws"
	"github.com/scottd018/policy-gen/internal/pkg/input"
	"github.com/scottd018/policy-gen/internal/pkg/processor"
)

const awsPolicyGenExample = `
# generate policies using sensible defaults
policy-gen aws

# generate policies from files located at input path ./input and write 
# discovered policies to output path ./output while forcefully overwriting
# any overlapping policies in the ./output directory.
policy-gen aws --input-path=./input --output-path=./output --force

# generate policies with debug logging
policy-gen aws --debug

# generate policies and associated documentation at ./output/README.md
policy-gen aws --output-path=./output --documentation=README.md
`

func NewCommand() *cobra.Command {
	flags := input.NewFlags()

	// create the command
	command := &cobra.Command{
		Use:     "aws",
		Short:   "Generate AWS IAM policies",
		Long:    `Generate AWS IAM policies`,
		RunE:    func(cmd *cobra.Command, args []string) error { return run(flags) },
		Example: awsPolicyGenExample,
	}

	// initialize the flags
	flags.Initialize(command)

	return command
}

func run(flags input.Flags) error {
	// convert our user input into a configuration for the processor
	config, err := flags.ToProcessorConfig()
	if err != nil {
		return fmt.Errorf("unable to convert flags into a processor config - %w", err)
	}

	// create the processor
	processor, err := processor.NewProcessor(
		config,
		aws.MarkerDefinition(),
		aws.Marker{},
		&aws.PolicyDocumentGenerator{Directory: config.OutputDirectory},
	)
	if err != nil {
		return err
	}

	// execute
	if err := processor.Process(); err != nil {
		return fmt.Errorf("unable to process markers - %w", err)
	}

	return nil
}
