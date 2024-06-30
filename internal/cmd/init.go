package cmd

import (
	"fmt"

	"github.com/adrianpk/gohermes/internal/handler"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the directory structure",
	Run: func(cmd *cobra.Command, args []string) {
		dirs := []string{
			"content",
			"content/root",
			"content/root/blog",
			"content/root/pages",
			"content/root/series",
			"content/section/blog",
			"content/section/pages",
			"content/section/series",
			"layout/default",
			"output",
			"store",
		}

		err := handler.InitDirs(dirs)
		if err != nil {
			fmt.Println("Error initializing directories:", err)
			return
		}

		fmt.Println("Directory structure initialized.")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
