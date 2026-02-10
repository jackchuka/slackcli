package users

import (
	"github.com/spf13/cobra"

	"github.com/jackchuka/slackcli/internal/cmdutil"
	"github.com/jackchuka/slackcli/internal/slack"
)

func NewUsersCmd() *cobra.Command {
	usersCmd := &cobra.Command{
		Use:   "users",
		Short: "Manage users",
	}
	usersCmd.AddCommand(newListCmd())
	usersCmd.AddCommand(newInfoCmd())
	usersCmd.AddCommand(newPresenceCmd())
	return usersCmd
}

func newListCmd() *cobra.Command {
	var limit int
	var all bool

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List users",
		RunE: func(c *cobra.Command, args []string) error {
			rc := cmdutil.GetRunContext(c.Context())
			result, err := rc.Client.ListUsers(slack.PaginationParams{
				Limit: limit,
				All:   all,
			})
			if err != nil {
				return err
			}
			return rc.Formatter.Format(result)
		},
	}
	listCmd.Flags().IntVar(&limit, "limit", 100, "Number of users per page")
	listCmd.Flags().BoolVar(&all, "all", false, "Fetch all users")
	return listCmd
}

func newInfoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "info <user-id>",
		Short: "Get user info",
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			rc := cmdutil.GetRunContext(c.Context())
			user, err := rc.Client.GetUserInfo(args[0])
			if err != nil {
				return err
			}
			return rc.Formatter.Format(user)
		},
	}
}

func newPresenceCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "presence <user-id>",
		Short: "Get user presence",
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			rc := cmdutil.GetRunContext(c.Context())
			presence, err := rc.Client.GetUserPresence(args[0])
			if err != nil {
				return err
			}
			return rc.Formatter.Format(map[string]string{
				"user_id":  args[0],
				"presence": presence,
			})
		},
	}
}
