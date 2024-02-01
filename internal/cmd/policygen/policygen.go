package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/scottd018/policy-gen/internal/cmd/policygen/aws"
	"github.com/scottd018/policy-gen/internal/cmd/policygen/version"
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(command *cobra.Command) {
	err := command.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// main executes the main program loop.
func main() {
	policygen := &cobra.Command{
		Use:   "policy-gen",
		Short: "Generate policies for public cloud providers",
		Long:  `Generate policies for public cloud providers`,
	}

	policygen.AddCommand(aws.NewCommand())
	policygen.AddCommand(version.NewCommand())
	Execute(policygen)
}
