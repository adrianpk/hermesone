package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hermes",
	Short: "Hermes is a CLI tool for managing your project",
}

func Execute() error {
	return rootCmd.Execute()
}
