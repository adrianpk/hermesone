package main

import (
	"embed"

	"github.com/adrianpk/gohermes/internal/cmd"
)

//go:embed layout/default/default.html
var layoutFS embed.FS

func main() {
	rootCmd := cmd.NewRootCmd()
	rootCmd.AddCommand(cmd.NewInitCmd(layoutFS))
	rootCmd.AddCommand(cmd.NewGenCmd())	
	rootCmd.Execute()
}
