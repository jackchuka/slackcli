package channels

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jackchuka/slackcli/internal/cmdutil"
	"github.com/jackchuka/slackcli/internal/slack"
)

func NewChannelsCmd() *cobra.Command {
	channelsCmd := &cobra.Command{
		Use:   "channels",
		Short: "Manage channels",
	}
	channelsCmd.AddCommand(newListCmd())
	channelsCmd.AddCommand(newInfoCmd())
	channelsCmd.AddCommand(newCreateCmd())
	channelsCmd.AddCommand(newArchiveCmd())
	channelsCmd.AddCommand(newInviteCmd())
	channelsCmd.AddCommand(newKickCmd())
	channelsCmd.AddCommand(newTopicCmd())
	channelsCmd.AddCommand(newPurposeCmd())
	return channelsCmd
}

func newListCmd() *cobra.Command {
	var cursor string
	var limit int
	var all bool

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List channels",
		RunE: func(c *cobra.Command, args []string) error {
			rc := cmdutil.GetRunContext(c.Context())
			result, err := rc.Client.ListChannels(slack.PaginationParams{
				Cursor: cursor,
				Limit:  limit,
				All:    all,
			})
			if err != nil {
				return err
			}
			return rc.Formatter.Format(result)
		},
	}
	listCmd.Flags().StringVar(&cursor, "cursor", "", "Pagination cursor")
	listCmd.Flags().IntVar(&limit, "limit", 100, "Number of channels per page")
	listCmd.Flags().BoolVar(&all, "all", false, "Fetch all channels (auto-paginate)")
	return listCmd
}

func newInfoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "info <channel-id>",
		Short: "Get channel info",
		Args:  cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			rc := cmdutil.GetRunContext(c.Context())
			ch, err := rc.Client.GetChannelInfo(args[0])
			if err != nil {
				return err
			}
			return rc.Formatter.Format(ch)
		},
	}
}

func newCreateCmd() *cobra.Command {
	var private bool

	createCmd := &cobra.Command{
		Use:         "create <name>",
		Short:       "Create a channel",
		Args:        cobra.ExactArgs(1),
		Annotations: map[string]string{"mode": "write"},
		RunE: func(c *cobra.Command, args []string) error {
			rc := cmdutil.GetRunContext(c.Context())
			ch, err := rc.Client.CreateChannel(args[0], private)
			if err != nil {
				return err
			}
			return rc.Formatter.Format(ch)
		},
	}
	createCmd.Flags().BoolVar(&private, "private", false, "Create a private channel")
	return createCmd
}

func newArchiveCmd() *cobra.Command {
	return &cobra.Command{
		Use:         "archive <channel-id>",
		Short:       "Archive a channel",
		Args:        cobra.ExactArgs(1),
		Annotations: map[string]string{"mode": "write"},
		RunE: func(c *cobra.Command, args []string) error {
			rc := cmdutil.GetRunContext(c.Context())
			if err := rc.Client.ArchiveChannel(args[0]); err != nil {
				return err
			}
			return rc.Formatter.Format(map[string]string{
				"status":  "archived",
				"channel": args[0],
			})
		},
	}
}

func newInviteCmd() *cobra.Command {
	return &cobra.Command{
		Use:         "invite <channel-id> <user-id>...",
		Short:       "Invite users to a channel",
		Args:        cobra.MinimumNArgs(2),
		Annotations: map[string]string{"mode": "write"},
		RunE: func(c *cobra.Command, args []string) error {
			rc := cmdutil.GetRunContext(c.Context())
			if err := rc.Client.InviteToChannel(args[0], args[1:]...); err != nil {
				return err
			}
			return rc.Formatter.Format(map[string]any{
				"status":  "invited",
				"channel": args[0],
				"users":   args[1:],
			})
		},
	}
}

func newKickCmd() *cobra.Command {
	return &cobra.Command{
		Use:         "kick <channel-id> <user-id>",
		Short:       "Remove a user from a channel",
		Args:        cobra.ExactArgs(2),
		Annotations: map[string]string{"mode": "write"},
		RunE: func(c *cobra.Command, args []string) error {
			rc := cmdutil.GetRunContext(c.Context())
			if err := rc.Client.KickFromChannel(args[0], args[1]); err != nil {
				return err
			}
			return rc.Formatter.Format(map[string]string{
				"status":  "removed",
				"channel": args[0],
				"user":    args[1],
			})
		},
	}
}

func newTopicCmd() *cobra.Command {
	return &cobra.Command{
		Use:         "topic <channel-id> <topic>",
		Short:       "Set channel topic",
		Args:        cobra.ExactArgs(2),
		Annotations: map[string]string{"mode": "write"},
		RunE: func(c *cobra.Command, args []string) error {
			rc := cmdutil.GetRunContext(c.Context())
			if err := rc.Client.SetChannelTopic(args[0], args[1]); err != nil {
				return err
			}
			return rc.Formatter.Format(map[string]string{
				"status":  "updated",
				"channel": args[0],
				"topic":   args[1],
			})
		},
	}
}

func newPurposeCmd() *cobra.Command {
	return &cobra.Command{
		Use:         "purpose <channel-id> <purpose>",
		Short:       "Set channel purpose",
		Args:        cobra.ExactArgs(2),
		Annotations: map[string]string{"mode": "write"},
		RunE: func(c *cobra.Command, args []string) error {
			rc := cmdutil.GetRunContext(c.Context())
			if err := rc.Client.SetChannelPurpose(args[0], args[1]); err != nil {
				return err
			}
			return rc.Formatter.Format(map[string]string{
				"status":  "updated",
				"channel": args[0],
				"purpose": args[1],
			})
		},
	}
}

func init() {
	_ = fmt.Sprintf // avoid unused import error
}
