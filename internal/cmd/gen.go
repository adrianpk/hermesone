package cmd

import (
	"log"

	"github.com/adrianpk/gohermes/internal/handler"
	"github.com/spf13/cobra"
)

func NewGenCmd() *cobra.Command {
	return genCmd
}

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate HTML from Markdown",
	Run: func(cmd *cobra.Command, args []string) {
		err := handler.GenHTML()
		if err != nil {
			log.Println("Error generating HTML:", err)
			return
		}

		log.Println("HTML generated.")
	},
}

