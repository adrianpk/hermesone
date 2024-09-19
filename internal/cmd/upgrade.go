package cmd

import (
	"embed"
	"github.com/spf13/cobra"
	"github.com/adrianpk/gohermes/internal/handler"
)

func NewUpgradeCmd(layoutFS embed.FS) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "upgrade the existing layout structure",
		RunE: func(cmd *cobra.Command, args []string) error {
			dirs := []string{
				"layout/default/articles",
				"layout/default/blog",
				"layout/default/pages",
				"layout/default/series",
			}
			return handler.Upgrade(dirs, layoutFS)
		},
	}

	return cmd
}
