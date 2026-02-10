package reactions

import (
	"github.com/spf13/cobra"

	"github.com/jackchuka/slackcli/internal/cmdutil"
	"github.com/jackchuka/slackcli/internal/slack"
)

func NewReactionsCmd() *cobra.Command {
	reactionsCmd := &cobra.Command{
		Use:   "reactions",
		Short: "Manage reactions",
	}
	reactionsCmd.AddCommand(newAddCmd())
	reactionsCmd.AddCommand(newRemoveCmd())
	reactionsCmd.AddCommand(newListCmd())
	return reactionsCmd
}

func newAddCmd() *cobra.Command {
	var channelID string
	var timestamp string
	var name string

	addCmd := &cobra.Command{
		Use:         "add",
		Short:       "Add a reaction",
		Annotations: map[string]string{"mode": "write"},
		RunE: func(c *cobra.Command, args []string) error {
			rc := cmdutil.GetRunContext(c.Context())
			if err := rc.Client.AddReaction(channelID, timestamp, name); err != nil {
				return err
			}
			return rc.Formatter.Format(map[string]string{
				"status":    "added",
				"reaction":  name,
				"channel":   channelID,
				"timestamp": timestamp,
			})
		},
	}
	addCmd.Flags().StringVar(&channelID, "channel", "", "Channel ID (required)")
	_ = addCmd.MarkFlagRequired("channel")
	addCmd.Flags().StringVar(&timestamp, "timestamp", "", "Message timestamp (required)")
	_ = addCmd.MarkFlagRequired("timestamp")
	addCmd.Flags().StringVar(&name, "name", "", "Reaction emoji name (required)")
	_ = addCmd.MarkFlagRequired("name")
	return addCmd
}

func newRemoveCmd() *cobra.Command {
	var channelID string
	var timestamp string
	var name string

	removeCmd := &cobra.Command{
		Use:         "remove",
		Short:       "Remove a reaction",
		Annotations: map[string]string{"mode": "write"},
		RunE: func(c *cobra.Command, args []string) error {
			rc := cmdutil.GetRunContext(c.Context())
			if err := rc.Client.RemoveReaction(channelID, timestamp, name); err != nil {
				return err
			}
			return rc.Formatter.Format(map[string]string{
				"status":    "removed",
				"reaction":  name,
				"channel":   channelID,
				"timestamp": timestamp,
			})
		},
	}
	removeCmd.Flags().StringVar(&channelID, "channel", "", "Channel ID (required)")
	_ = removeCmd.MarkFlagRequired("channel")
	removeCmd.Flags().StringVar(&timestamp, "timestamp", "", "Message timestamp (required)")
	_ = removeCmd.MarkFlagRequired("timestamp")
	removeCmd.Flags().StringVar(&name, "name", "", "Reaction emoji name (required)")
	_ = removeCmd.MarkFlagRequired("name")
	return removeCmd
}

func newListCmd() *cobra.Command {
	var userID string
	var limit int
	var all bool

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List reactions for a user",
		RunE: func(c *cobra.Command, args []string) error {
			rc := cmdutil.GetRunContext(c.Context())
			result, err := rc.Client.ListReactions(userID, slack.PaginationParams{
				Limit: limit,
				All:   all,
			})
			if err != nil {
				return err
			}
			return rc.Formatter.Format(result)
		},
	}
	listCmd.Flags().StringVar(&userID, "user", "", "User ID (defaults to authenticated user)")
	listCmd.Flags().IntVar(&limit, "limit", 100, "Number of reactions per page")
	listCmd.Flags().BoolVar(&all, "all", false, "Fetch all reactions (auto-paginate)")
	return listCmd
}
