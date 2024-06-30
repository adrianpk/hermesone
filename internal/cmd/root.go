package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd *cobra.Command

func NewRootCmd() *cobra.Command {
	if rootCmd == nil {
		rootCmd = &cobra.Command{
			Use:   "hermes",
			Short: "Hermes is a CLI tool for managing your SSG projec",
		}
	}

	return rootCmd
}

func Execute() error {
	return rootCmd.Execute()
}
