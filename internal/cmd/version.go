package cmd

import (
	"github.com/spf13/cobra"

	"github.com/jackchuka/slackcli/internal/cmdutil"
	"github.com/jackchuka/slackcli/internal/version"
)

func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		RunE: func(cmd *cobra.Command, args []string) error {
			rc := cmdutil.GetRunContext(cmd.Context())
			return rc.Formatter.Format(map[string]string{
				"version":    version.Version,
				"commit":     version.Commit,
				"build_date": version.BuildDate,
			})
		},
	}
}
