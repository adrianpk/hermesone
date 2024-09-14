package cmd

import (
	"embed"
	"log"

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
				"content/root/pages",
				"content/root/articles",
				"content/root/blog",
				"content/root/series",
				"content/section/pages",
				"content/section/articles",
				"content/section/blog",
				"content/section/series",
				"layout/default",
				"layout/default/pages",
				"layout/default/articles",
				"layout/default/blog",
				"layout/default/series",
				"output",
				"store",
			}

			err := handler.InitDirs(dirs, layoutFS)
			if err != nil {
				log.Println("Error initializing directories:", err)
				return
			}

			log.Println("Directory structure initialized.")
		},
	}
}

