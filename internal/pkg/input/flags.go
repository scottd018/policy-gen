package input

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/scottd018/policy-gen/internal/pkg/files"
	"github.com/scottd018/policy-gen/internal/pkg/processor"
)

// FlagInput represents an individual flag value as determined from user input.
type FlagInput struct {
	CommandFunc    func(*cobra.Command, *FlagInput)
	StringDefault  string
	StringValue    string
	BooleanDefault bool
	BooleanValue   bool
	Description    string
	Short          string
	Required       bool
}

// Flags represents the user input into the command line.
type Flags map[string]*FlagInput

// NewFlags returns a new set of flags for the policy-gen CLI.
func NewFlags() Flags {
	return Flags{
		FlagInputPath: &FlagInput{
			StringDefault: FlagInputPathDefault,
			Description:   FlagInputPathDescription,
			Short:         FlagInputPathShort,
			Required:      true,
			CommandFunc: func(command *cobra.Command, input *FlagInput) {
				command.Flags().StringVarP(&input.StringValue, FlagInputPath, input.Short, input.StringDefault, input.Description)
			},
		},
		FlagOutputPath: &FlagInput{
			StringDefault: FlagOutputPathDefault,
			Description:   FlagOutputPathDescription,
			Short:         FlagOutputPathShort,
			Required:      true,
			CommandFunc: func(command *cobra.Command, input *FlagInput) {
				command.Flags().StringVarP(&input.StringValue, FlagOutputPath, input.Short, input.StringDefault, input.Description)
			},
		},
		FlagDocumentation: &FlagInput{
			StringDefault: FlagDocumentationDefault,
			Description:   FlagDocumentationDescription,
			Short:         FlagDocumentationShort,
			Required:      false,
			CommandFunc: func(command *cobra.Command, input *FlagInput) {
				command.Flags().StringVarP(&input.StringValue, FlagDocumentation, input.Short, input.StringDefault, input.Description)
			},
		},
		FlagForce: &FlagInput{
			BooleanDefault: FlagForceDefault,
			Description:    FlagForceDescription,
			Short:          FlagForceShort,
			Required:       false,
			CommandFunc: func(command *cobra.Command, input *FlagInput) {
				command.Flags().BoolVarP(&input.BooleanValue, FlagForce, input.Short, input.BooleanDefault, input.Description)
			},
		},
		FlagRecursive: &FlagInput{
			BooleanDefault: FlagRecursiveDefault,
			Description:    FlagRecursiveDescription,
			Short:          FlagRecursiveShort,
			Required:       false,
			CommandFunc: func(command *cobra.Command, input *FlagInput) {
				command.Flags().BoolVarP(&input.BooleanValue, FlagRecursive, input.Short, input.BooleanDefault, input.Description)
			},
		},
		FlagDebug: &FlagInput{
			BooleanDefault: FlagDebugDefault,
			Description:    FlagDebugDescription,
			Required:       false,
			CommandFunc: func(command *cobra.Command, input *FlagInput) {
				command.Flags().BoolVar(&input.BooleanValue, FlagDebug, input.BooleanDefault, input.Description)
			},
		},
	}
}

// Initialize initializes a set of flags by running adding the flags to the command using the CommandFunc.
func (flags Flags) Initialize(command *cobra.Command) {
	for flag, input := range flags {
		input.CommandFunc(command, flags.For(flag))
	}
}

// ToProcessorConfig processes the raw input flags validates them, and converts them to an processor configuration.
func (flags Flags) ToProcessorConfig() (*processor.Config, error) {
	// ensure required string values have values set
	for flag, input := range flags {
		if input.Required && input.StringValue == "" {
			return nil, fmt.Errorf("missing value for required flag: [--%s]", flag)
		}
	}

	// validate existence of directory objects and add them to the processor
	inputDirectory, err := files.NewDirectory(flags.For(FlagInputPath).StringValue, files.WithPreExistingDirectory)
	if err != nil {
		return nil, fmt.Errorf("invalid flag: [--%s] - %w", FlagInputPath, err)
	}

	outputDirectory, err := files.NewDirectory(flags.For(FlagOutputPath).StringValue, files.WithPreExistingDirectory)
	if err != nil {
		return nil, fmt.Errorf("invalid flag: [--%s] - %w", FlagOutputPath, err)
	}

	// validate existence of file objects and add them to the processor
	var documentationFile *files.File

	documentationInput := flags.For(FlagDocumentation).StringValue
	if documentationInput != "" {
		documentationFile, err = files.NewFile(documentationInput, files.WithPreExistingDirectory)
		if err != nil {
			return nil, fmt.Errorf("invalid flag: [--%s] - %w", FlagDocumentation, err)
		}
	}

	return &processor.Config{
		InputDirectory:    inputDirectory,
		OutputDirectory:   outputDirectory,
		DocumentationFile: documentationFile,
		Force:             flags.For(FlagForce).BooleanValue,
		Debug:             flags.For(FlagDebug).BooleanValue,
	}, nil
}

// For returns the FlagInput for a particular flag.
func (flags Flags) For(flag string) *FlagInput {
	return flags[flag]
}
