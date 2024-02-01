package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "unstable"

const versionExample = `
policy-gen version
`

func NewCommand() *cobra.Command {
	// create the command
	command := &cobra.Command{
		Use:     "version",
		Short:   "Print version",
		Long:    `Print version`,
		Run:     func(cmd *cobra.Command, args []string) { run() },
		Example: versionExample,
	}

	return command
}

//nolint:forbidigo
func run() {
	fmt.Printf("%s\n", version)
}
