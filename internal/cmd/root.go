package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/jackchuka/slackcli/internal/auth"
	"github.com/jackchuka/slackcli/internal/cmdutil"
	"github.com/jackchuka/slackcli/internal/config"
	"github.com/jackchuka/slackcli/internal/output"
	"github.com/jackchuka/slackcli/internal/slack"

	authcmd "github.com/jackchuka/slackcli/internal/cmd/auth"
	channelscmd "github.com/jackchuka/slackcli/internal/cmd/channels"
	filescmd "github.com/jackchuka/slackcli/internal/cmd/files"
	mcpcmd "github.com/jackchuka/slackcli/internal/cmd/mcp"
	messagescmd "github.com/jackchuka/slackcli/internal/cmd/messages"
	reactionscmd "github.com/jackchuka/slackcli/internal/cmd/reactions"
	userscmd "github.com/jackchuka/slackcli/internal/cmd/users"
)

var (
	flagToken     string
	flagWorkspace string
	flagOutput    string
	flagReadOnly  bool
)

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:           "slackcli",
		Short:         "Slack CLI & MCP Server",
		Long:          "A CLI tool and MCP server for programmatic Slack access, optimized for LLM/AI agents.",
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Skip client initialization for auth and version commands
			if cmd.Name() == "version" || cmd.Parent().Name() == "auth" || cmd.Name() == "serve" {
				if err := initContext(cmd, false); err != nil {
					return err
				}
			} else {
				if err := initContext(cmd, true); err != nil {
					return err
				}
			}

			// Enforce read-only mode
			if flagReadOnly {
				if cmd.Annotations != nil && cmd.Annotations["mode"] == "write" {
					return fmt.Errorf("command %q is a write operation and cannot be used in read-only mode", cmd.CommandPath())
				}
			}
			return nil
		},
	}

	rootCmd.PersistentFlags().StringVar(&flagToken, "token", "", "Slack API token")
	rootCmd.PersistentFlags().StringVarP(&flagWorkspace, "workspace", "w", "", "Workspace name")
	rootCmd.PersistentFlags().StringVarP(&flagOutput, "output", "o", "", "Output format (json|table)")
	rootCmd.PersistentFlags().BoolVar(&flagReadOnly, "read-only", false, "Restrict to read-only operations (reject writes)")

	rootCmd.AddCommand(NewVersionCmd())
	rootCmd.AddCommand(authcmd.NewAuthCmd())
	rootCmd.AddCommand(channelscmd.NewChannelsCmd())
	rootCmd.AddCommand(messagescmd.NewMessagesCmd())
	rootCmd.AddCommand(userscmd.NewUsersCmd())
	rootCmd.AddCommand(reactionscmd.NewReactionsCmd())
	rootCmd.AddCommand(filescmd.NewFilesCmd())
	rootCmd.AddCommand(mcpcmd.NewMCPCmd())

	return rootCmd
}

func initContext(cmd *cobra.Command, needsClient bool) error {
	cfg, err := config.Load("")
	if err != nil {
		return err
	}

	writers := output.DefaultWriters()

	var formatter output.Formatter
	switch flagOutput {
	case "table":
		formatter = output.NewTableFormatter(writers.Out)
	case "json":
		formatter = output.NewJSONFormatter(writers.Out)
	default:
		if output.IsTTY(os.Stdout) {
			formatter = output.NewTableFormatter(writers.Out)
		} else {
			formatter = output.NewJSONFormatter(writers.Out)
		}
	}

	resolver := auth.NewResolver(flagToken, cfg.ActiveToken)

	rc := &cmdutil.RunContext{
		Config:    cfg,
		Formatter: formatter,
		Writers:   writers,
		Resolver:  resolver,
		ReadOnly:  flagReadOnly,
	}

	if needsClient {
		token := resolver.Resolve()
		if token == "" {
			output.PrintError(writers, "no token found. Run 'slackcli auth login' or set SLACK_TOKEN")
			os.Exit(2)
		}
		rc.Client = slack.NewClient(token)
	}

	cmd.SetContext(cmdutil.SetRunContext(cmd.Context(), rc))
	return nil
}

func Execute() error {
	return NewRootCmd().Execute()
}
