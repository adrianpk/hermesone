package main

import (
	"embed"

	"github.com/adrianpk/gohermes/internal/cmd"
)

const ver = "0.0.1"

//go:embed layout/default
var layoutFS embed.FS

func main() {
	rootCmd := cmd.NewRootCmd(ver)
	rootCmd.AddCommand(cmd.NewInitCmd(layoutFS))
	rootCmd.AddCommand(cmd.NewGenCmd())	
	rootCmd.Execute()
}
