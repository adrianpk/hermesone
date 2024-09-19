package cmd

import (
	"embed"
	"log"

	"github.com/adrianpk/gohermes/internal/handler"
	"github.com/adrianpk/gohermes/internal/hermes"
	"github.com/spf13/cobra"
)

const defName = "hermes-site"

func NewInitCmd(layoutFS embed.FS) *cobra.Command {
	var projectName string

	cmd := &cobra.Command{
		Use:   "init",
		Short: "initialize the directory structure",
		Run: func(cmd *cobra.Command, args []string) {
			if projectName == "" {
				projectName = defName
			}

			dirs := []string{
				"content",
				"content/root",
				"content/root/page",
				"content/root/article",
				"content/root/blog",
				"content/root/series",
				"content/section/page",
				"content/section/article",
				"content/section/blog",
				"content/section/series",
				"layout/default",
				"layout/default/page",
				"layout/default/article",
				"layout/default/blog",
				"layout/default/series",
				"output",
				"store",
			}

			err := handler.InitDirs(dirs, layoutFS)
			if err != nil {
				log.Println("error initializing directories:", err)
				return
			}

			err = hermes.NewCfgFile(projectName)
			if err != nil {
				log.Println("error creating hermes.yml:", err)
				return
			}

			log.Println("directory structure initialized.")
		},
	}

	cmd.Flags().StringVarP(&projectName, "name", "n", "", "name of the project")
	return cmd
}
