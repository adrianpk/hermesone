package cmd

import (
	"embed"
	"fmt"

	"github.com/adrianpk/gohermes/internal/handler"
	"github.com/spf13/cobra"
)

func NewInitCmd(layoutFS embed.FS) *cobra.Command {
	return &cobra.Command{
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
				"layout/default/blog",
				"layout/default/pages",
				"layout/default/series",
				"output",
				"store",
			}

			err := handler.InitDirs(dirs, layoutFS)
			if err != nil {
				fmt.Println("Error initializing directories:", err)
				return
			}

			fmt.Println("Directory structure initialized.")
		},
	}
}

