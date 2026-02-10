package mcp

import (
	"fmt"

	"github.com/mark3labs/mcp-go/server"
	"github.com/spf13/cobra"

	"github.com/jackchuka/slackcli/internal/cmdutil"
	mcpserver "github.com/jackchuka/slackcli/internal/mcp"
	"github.com/jackchuka/slackcli/internal/slack"
)

func NewMCPCmd() *cobra.Command {
	mcpCmd := &cobra.Command{
		Use:   "mcp",
		Short: "MCP server commands",
	}
	mcpCmd.AddCommand(newServeCmd())
	return mcpCmd
}

func newServeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Start MCP server (stdio transport)",
		RunE: func(c *cobra.Command, args []string) error {
			rc := cmdutil.GetRunContext(c.Context())
			token := rc.Resolver.Resolve()
			if token == "" {
				return fmt.Errorf("no token found. Set SLACK_TOKEN or run 'slackcli auth login'")
			}
			client := slack.NewClient(token)
			s := mcpserver.NewServer(client, rc.ReadOnly)
			return server.ServeStdio(s)
		},
	}
}
