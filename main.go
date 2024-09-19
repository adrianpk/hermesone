package main

import (
	"embed"
	"fmt"
	"os"

	"github.com/adrianpk/gohermes/internal/cmd"
)

const ver = "0.0.1"

//go:embed layout
var layoutFS embed.FS

func main() {
	rootCmd := cmd.NewRootCmd(ver)
	rootCmd.AddCommand(cmd.NewInitCmd(layoutFS))
	rootCmd.AddCommand(cmd.NewGenCmd())
	rootCmd.AddCommand(cmd.NewUpgradeCmd(layoutFS))
	rootCmd.AddCommand(cmd.NewNewCmd())

	if len(os.Args) > 1 && os.Args[1] == "help" {
		rootCmd.Usage()
		return
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
