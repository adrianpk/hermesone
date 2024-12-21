package cmd

import (
	"embed"

	"github.com/adrianpk/gohermes/internal/handler"
	"github.com/spf13/cobra"
)

func NewUpgradeCmd(layoutFS embed.FS) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "upgrade the existing layout structure",
		RunE: func(cmd *cobra.Command, args []string) error {
			dirs := []string{
				"assets/layout/default/article",
				"assets/layout/default/blog",
				"assets/layout/default/page",
				"assets/layout/default/series",
			}
			return handler.Upgrade(dirs, layoutFS)
		},
	}

	return cmd
}
