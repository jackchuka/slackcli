package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/jackchuka/slackcli/internal/slack"
)

func registerMessageTools(s *server.MCPServer, client slack.Service, readOnly bool) {
	s.AddTool(mcp.NewTool("list_messages",
		mcp.WithDescription("List messages in a Slack channel"),
		mcp.WithString("channel_id", mcp.Required(), mcp.Description("Channel ID")),
		mcp.WithNumber("limit", mcp.Description("Max messages to return"), mcp.DefaultNumber(100)),
		mcp.WithBoolean("all", mcp.Description("Fetch all messages")),
		mcp.WithString("cursor", mcp.Description("Pagination cursor")),
	), makeListMessages(client))

	if readOnly {
		return
	}

	s.AddTool(mcp.NewTool("send_message",
		mcp.WithDescription("Send a message to a Slack channel"),
		mcp.WithString("channel_id", mcp.Required(), mcp.Description("Channel ID")),
		mcp.WithString("text", mcp.Required(), mcp.Description("Message text")),
		mcp.WithString("thread_ts", mcp.Description("Thread timestamp for replies")),
	), makeSendMessage(client))

	s.AddTool(mcp.NewTool("edit_message",
		mcp.WithDescription("Edit an existing message"),
		mcp.WithString("channel_id", mcp.Required(), mcp.Description("Channel ID")),
		mcp.WithString("timestamp", mcp.Required(), mcp.Description("Message timestamp")),
		mcp.WithString("text", mcp.Required(), mcp.Description("New message text")),
	), makeEditMessage(client))

	s.AddTool(mcp.NewTool("delete_message",
		mcp.WithDescription("Delete a message"),
		mcp.WithString("channel_id", mcp.Required(), mcp.Description("Channel ID")),
		mcp.WithString("timestamp", mcp.Required(), mcp.Description("Message timestamp")),
	), makeDeleteMessage(client))
}

func makeListMessages(client slack.Service) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		channelID, err := request.RequireString("channel_id")
		if err != nil {
			return errResult(err), nil
		}
		limit := request.GetInt("limit", 100)
		all := request.GetBool("all", false)
		cursor := request.GetString("cursor", "")

		result, err := client.ListMessages(slack.ListMessagesParams{
			ChannelID:  channelID,
			Pagination: slack.PaginationParams{Cursor: cursor, Limit: limit, All: all},
		})
		if err != nil {
			return errResult(err), nil
		}
		return mcp.NewToolResultText(toJSON(result)), nil
	}
}

func makeSendMessage(client slack.Service) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		channelID, err := request.RequireString("channel_id")
		if err != nil {
			return errResult(err), nil
		}
		text, err := request.RequireString("text")
		if err != nil {
			return errResult(err), nil
		}
		threadTS := request.GetString("thread_ts", "")

		msg, err := client.SendMessage(slack.SendMessageParams{
			ChannelID: channelID,
			Text:      text,
			ThreadTS:  threadTS,
		})
		if err != nil {
			return errResult(err), nil
		}
		return mcp.NewToolResultText(toJSON(msg)), nil
	}
}

func makeEditMessage(client slack.Service) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		channelID, err := request.RequireString("channel_id")
		if err != nil {
			return errResult(err), nil
		}
		timestamp, err := request.RequireString("timestamp")
		if err != nil {
			return errResult(err), nil
		}
		text, err := request.RequireString("text")
		if err != nil {
			return errResult(err), nil
		}

		msg, err := client.EditMessage(channelID, timestamp, text)
		if err != nil {
			return errResult(err), nil
		}
		return mcp.NewToolResultText(toJSON(msg)), nil
	}
}

func makeDeleteMessage(client slack.Service) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		channelID, err := request.RequireString("channel_id")
		if err != nil {
			return errResult(err), nil
		}
		timestamp, err := request.RequireString("timestamp")
		if err != nil {
			return errResult(err), nil
		}
		if err := client.DeleteMessage(channelID, timestamp); err != nil {
			return errResult(err), nil
		}
		return mcp.NewToolResultText(toJSON(map[string]string{"status": "deleted", "channel_id": channelID, "timestamp": timestamp})), nil
	}
}
