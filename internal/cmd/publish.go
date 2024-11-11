package cmd

import (
	"log"

	"github.com/adrianpk/gohermes/internal/handler"
	"github.com/spf13/cobra"
)

func NewPublishCmd() *cobra.Command {
	return publishCmd
}

var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish the generated content to GitHub Pages",
	Run: func(cmd *cobra.Command, args []string) {
		err := handler.Publish()
		if err != nil {
			log.Println("error publishing:", err)
			return
		}

		log.Println("Successfully published to GitHub Pages")
	},
}
