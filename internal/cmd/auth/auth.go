package auth

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jackchuka/slackcli/internal/cmdutil"
	"github.com/jackchuka/slackcli/internal/config"
	"github.com/jackchuka/slackcli/internal/slack"
)

func NewAuthCmd() *cobra.Command {
	authCmd := &cobra.Command{
		Use:   "auth",
		Short: "Manage authentication",
	}
	authCmd.AddCommand(newLoginCmd())
	authCmd.AddCommand(newLogoutCmd())
	authCmd.AddCommand(newStatusCmd())
	authCmd.AddCommand(newListCmd())
	authCmd.AddCommand(newSwitchCmd())
	return authCmd
}

func newLoginCmd() *cobra.Command {
	var token string
	var name string

	loginCmd := &cobra.Command{
		Use:   "login",
		Short: "Login to a Slack workspace",
		RunE: func(c *cobra.Command, args []string) error {
			rc := cmdutil.GetRunContext(c.Context())

			// Verify token
			client := slack.NewClient(token)
			authResult, err := client.AuthTest()
			if err != nil {
				return fmt.Errorf("authentication failed: %w", err)
			}

			wsName := name
			if wsName == "" {
				wsName = authResult.Team
			}

			rc.Config.SetWorkspace(wsName, config.Workspace{
				Name:   wsName,
				Token:  token,
				TeamID: authResult.TeamID,
			})

			if err := rc.Config.Save(); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}

			return rc.Formatter.Format(map[string]string{
				"status":    "authenticated",
				"team":      authResult.Team,
				"user":      authResult.User,
				"team_id":   authResult.TeamID,
				"workspace": wsName,
			})
		},
	}
	loginCmd.Flags().StringVar(&token, "token", "", "Slack API token (xoxb-* or xoxp-*)")
	_ = loginCmd.MarkFlagRequired("token")
	loginCmd.Flags().StringVar(&name, "name", "", "Workspace name (defaults to team name)")
	return loginCmd
}

func newLogoutCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "logout [workspace]",
		Short: "Remove a workspace",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			rc := cmdutil.GetRunContext(c.Context())
			name := rc.Config.ActiveWorkspace
			if len(args) > 0 {
				name = args[0]
			}
			if _, ok := rc.Config.Workspaces[name]; !ok {
				return fmt.Errorf("workspace %q not found", name)
			}
			rc.Config.RemoveWorkspace(name)
			if err := rc.Config.Save(); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}
			return rc.Formatter.Format(map[string]string{
				"status":    "logged_out",
				"workspace": name,
			})
		},
	}
}

func newStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show current authentication status",
		RunE: func(c *cobra.Command, args []string) error {
			rc := cmdutil.GetRunContext(c.Context())
			token := rc.Resolver.Resolve()
			if token == "" {
				return rc.Formatter.Format(map[string]string{
					"status": "not_authenticated",
				})
			}
			client := slack.NewClient(token)
			result, err := client.AuthTest()
			if err != nil {
				return rc.Formatter.Format(map[string]any{
					"status": "error",
					"error":  err.Error(),
				})
			}
			return rc.Formatter.Format(map[string]any{
				"status":    "authenticated",
				"workspace": rc.Config.ActiveWorkspace,
				"team":      result.Team,
				"user":      result.User,
				"team_id":   result.TeamID,
			})
		},
	}
}

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List configured workspaces",
		RunE: func(c *cobra.Command, args []string) error {
			rc := cmdutil.GetRunContext(c.Context())
			type workspace struct {
				Name   string `json:"name"`
				TeamID string `json:"team_id"`
				Active bool   `json:"active"`
			}
			var workspaces []workspace
			for name, ws := range rc.Config.Workspaces {
				workspaces = append(workspaces, workspace{
					Name:   name,
					TeamID: ws.TeamID,
					Active: name == rc.Config.ActiveWorkspace,
				})
			}
			return rc.Formatter.Format(workspaces)
		},
	}
}

func newSwitchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "switch <workspace>",
		Short: "Switch active workspace",
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			rc := cmdutil.GetRunContext(c.Context())
			name := args[0]
			if _, ok := rc.Config.Workspaces[name]; !ok {
				return fmt.Errorf("workspace %q not found", name)
			}
			rc.Config.ActiveWorkspace = name
			if err := rc.Config.Save(); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}
			return rc.Formatter.Format(map[string]string{
				"status":    "switched",
				"workspace": name,
			})
		},
	}
}
