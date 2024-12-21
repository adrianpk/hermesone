package main

import (
	"embed"
	"fmt"
	"os"

	"github.com/adrianpk/gohermes/internal/cmd"
)

const ver = "0.0.1"

//go:embed assets/*
//go:embed assets/layout/default/partial/_index.html
var assetsFS embed.FS

func main() {
	rootCmd := cmd.NewRootCmd(ver)
	rootCmd.AddCommand(cmd.NewInitCmd(assetsFS))
	rootCmd.AddCommand(cmd.NewGenCmd())
	rootCmd.AddCommand(cmd.NewUpgradeCmd(assetsFS))
	rootCmd.AddCommand(cmd.NewNewCmd())
	rootCmd.AddCommand(cmd.NewPublishCmd())
	rootCmd.AddCommand(cmd.NewBackupCmd())

	if len(os.Args) > 1 && os.Args[1] == "help" {
		rootCmd.Usage()
		return
	}

	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
