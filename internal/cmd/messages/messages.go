package messages

import (
	"github.com/spf13/cobra"

	"github.com/jackchuka/slackcli/internal/cmdutil"
	"github.com/jackchuka/slackcli/internal/slack"
)

func NewMessagesCmd() *cobra.Command {
	messagesCmd := &cobra.Command{
		Use:   "messages",
		Short: "Manage messages",
	}
	messagesCmd.AddCommand(newListCmd())
	messagesCmd.AddCommand(newSendCmd())
	messagesCmd.AddCommand(newReplyCmd())
	messagesCmd.AddCommand(newEditCmd())
	messagesCmd.AddCommand(newDeleteCmd())
	messagesCmd.AddCommand(newSearchCmd())
	return messagesCmd
}

func newListCmd() *cobra.Command {
	var channelID string
	var cursor string
	var limit int
	var all bool

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List messages in a channel",
		RunE: func(c *cobra.Command, args []string) error {
			rc := cmdutil.GetRunContext(c.Context())
			result, err := rc.Client.ListMessages(slack.ListMessagesParams{
				ChannelID:  channelID,
				Pagination: slack.PaginationParams{Cursor: cursor, Limit: limit, All: all},
			})
			if err != nil {
				return err
			}
			return rc.Formatter.Format(result)
		},
	}
	listCmd.Flags().StringVar(&channelID, "channel", "", "Channel ID (required)")
	_ = listCmd.MarkFlagRequired("channel")
	listCmd.Flags().StringVar(&cursor, "cursor", "", "Pagination cursor")
	listCmd.Flags().IntVar(&limit, "limit", 100, "Number of messages per page")
	listCmd.Flags().BoolVar(&all, "all", false, "Fetch all messages")
	return listCmd
}

func newSendCmd() *cobra.Command {
	var channelID string
	var text string

	sendCmd := &cobra.Command{
		Use:         "send",
		Short:       "Send a message",
		Annotations: map[string]string{"mode": "write"},
		RunE: func(c *cobra.Command, args []string) error {
			rc := cmdutil.GetRunContext(c.Context())
			msg, err := rc.Client.SendMessage(slack.SendMessageParams{
				ChannelID: channelID,
				Text:      text,
			})
			if err != nil {
				return err
			}
			return rc.Formatter.Format(msg)
		},
	}
	sendCmd.Flags().StringVar(&channelID, "channel", "", "Channel ID (required)")
	_ = sendCmd.MarkFlagRequired("channel")
	sendCmd.Flags().StringVar(&text, "text", "", "Message text (required)")
	_ = sendCmd.MarkFlagRequired("text")
	return sendCmd
}

func newReplyCmd() *cobra.Command {
	var channelID string
	var threadTS string
	var text string

	replyCmd := &cobra.Command{
		Use:         "reply",
		Short:       "Reply to a thread",
		Annotations: map[string]string{"mode": "write"},
		RunE: func(c *cobra.Command, args []string) error {
			rc := cmdutil.GetRunContext(c.Context())
			msg, err := rc.Client.SendMessage(slack.SendMessageParams{
				ChannelID: channelID,
				Text:      text,
				ThreadTS:  threadTS,
			})
			if err != nil {
				return err
			}
			return rc.Formatter.Format(msg)
		},
	}
	replyCmd.Flags().StringVar(&channelID, "channel", "", "Channel ID (required)")
	_ = replyCmd.MarkFlagRequired("channel")
	replyCmd.Flags().StringVar(&threadTS, "thread-ts", "", "Thread timestamp (required)")
	_ = replyCmd.MarkFlagRequired("thread-ts")
	replyCmd.Flags().StringVar(&text, "text", "", "Message text (required)")
	_ = replyCmd.MarkFlagRequired("text")
	return replyCmd
}

func newEditCmd() *cobra.Command {
	var channelID string
	var timestamp string
	var text string

	editCmd := &cobra.Command{
		Use:         "edit",
		Short:       "Edit a message",
		Annotations: map[string]string{"mode": "write"},
		RunE: func(c *cobra.Command, args []string) error {
			rc := cmdutil.GetRunContext(c.Context())
			msg, err := rc.Client.EditMessage(channelID, timestamp, text)
			if err != nil {
				return err
			}
			return rc.Formatter.Format(msg)
		},
	}
	editCmd.Flags().StringVar(&channelID, "channel", "", "Channel ID (required)")
	_ = editCmd.MarkFlagRequired("channel")
	editCmd.Flags().StringVar(&timestamp, "timestamp", "", "Message timestamp (required)")
	_ = editCmd.MarkFlagRequired("timestamp")
	editCmd.Flags().StringVar(&text, "text", "", "New message text (required)")
	_ = editCmd.MarkFlagRequired("text")
	return editCmd
}

func newDeleteCmd() *cobra.Command {
	var channelID string
	var timestamp string

	deleteCmd := &cobra.Command{
		Use:         "delete",
		Short:       "Delete a message",
		Annotations: map[string]string{"mode": "write"},
		RunE: func(c *cobra.Command, args []string) error {
			rc := cmdutil.GetRunContext(c.Context())
			if err := rc.Client.DeleteMessage(channelID, timestamp); err != nil {
				return err
			}
			return rc.Formatter.Format(map[string]string{
				"status":    "deleted",
				"channel":   channelID,
				"timestamp": timestamp,
			})
		},
	}
	deleteCmd.Flags().StringVar(&channelID, "channel", "", "Channel ID (required)")
	_ = deleteCmd.MarkFlagRequired("channel")
	deleteCmd.Flags().StringVar(&timestamp, "timestamp", "", "Message timestamp (required)")
	_ = deleteCmd.MarkFlagRequired("timestamp")
	return deleteCmd
}

func newSearchCmd() *cobra.Command {
	var query string
	var sort string
	var sortDir string
	var limit int

	searchCmd := &cobra.Command{
		Use:   "search",
		Short: "Search messages",
		RunE: func(c *cobra.Command, args []string) error {
			rc := cmdutil.GetRunContext(c.Context())
			result, err := rc.Client.SearchMessages(slack.SearchParams{
				Query:      query,
				Sort:       sort,
				SortDir:    sortDir,
				Pagination: slack.PaginationParams{Limit: limit},
			})
			if err != nil {
				return err
			}
			return rc.Formatter.Format(result)
		},
	}
	searchCmd.Flags().StringVar(&query, "query", "", "Search query (required)")
	_ = searchCmd.MarkFlagRequired("query")
	searchCmd.Flags().StringVar(&sort, "sort", "timestamp", "Sort field (timestamp|score)")
	searchCmd.Flags().StringVar(&sortDir, "sort-dir", "desc", "Sort direction (asc|desc)")
	searchCmd.Flags().IntVar(&limit, "limit", 20, "Number of results")
	return searchCmd
}
