package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/scottd018/policy-gen/internal/cmd/policygen/aws"
)

// policygen represents the base command entrypoint.
var policygen = &cobra.Command{
	Use:   "policygen",
	Short: "Generate policies for public cloud providers",
	Long:  `Generate policies for public cloud providers`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := policygen.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// init adds all of the available subcommands to the root policygen command.
func init() {
	policygen.AddCommand(aws.NewCommand())
}

// main executes the main program loop.
func main() {
	Execute()
}
