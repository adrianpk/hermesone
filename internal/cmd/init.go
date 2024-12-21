package cmd

import (
	"embed"
	"log"

	"github.com/adrianpk/gohermes/internal/handler"
	"github.com/adrianpk/gohermes/internal/hermes"
	"github.com/spf13/cobra"
)

const defName = "hermes-site"

func NewInitCmd(assetsFS embed.FS) *cobra.Command {
	var projectName string
	var githubUser string

	cmd := &cobra.Command{
		Use:   "init",
		Short: "initialize the directory structure",
		Run: func(cmd *cobra.Command, args []string) {
			if projectName == "" {
				projectName = defName
			}

			if githubUser == "" {
				log.Println("GitHub username is required")
				return
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
				"output",
				"store",
			}

			err := handler.InitDirs(dirs, assetsFS)
			if err != nil {
				log.Println("error initializing directories:", err)
				return
			}

			err = hermes.NewCfgFile(projectName, githubUser)
			if err != nil {
				log.Println("error creating hermes.yml:", err)
				return
			}

			log.Println("directory structure initialized.")
		},
	}

	cmd.Flags().StringVarP(&projectName, "name", "n", "", "name of the project")
	cmd.Flags().StringVarP(&githubUser, "user", "u", "", "GitHub username")
	return cmd
}
