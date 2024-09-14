package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var rootCmd *cobra.Command

func NewRootCmd(ver string) *cobra.Command {
	log.Println("hermes version:", ver)

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
