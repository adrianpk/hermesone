package cmd

import (
	"log"

	"github.com/adrianpk/gohermes/internal/handler"
	"github.com/spf13/cobra"
)

func NewBackupCmd() *cobra.Command {
	return backupCmd
}

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup the Hermes project to GitHub",
	Run: func(cmd *cobra.Command, args []string) {
		err := handler.Backup()
		if err != nil {
			log.Println("error backing up:", err)
			return
		}

		log.Println("Successfully backed up to GitHub")
	},
}
