package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/jackchuka/slackcli/internal/slack"
)

func registerReactionTools(s *server.MCPServer, client slack.Service, readOnly bool) {
	s.AddTool(mcp.NewTool("list_reactions",
		mcp.WithDescription("List reactions for a user"),
		mcp.WithString("user_id", mcp.Description("User ID (defaults to authenticated user)")),
		mcp.WithNumber("limit", mcp.Description("Max items to return"), mcp.DefaultNumber(100)),
	), makeListReactions(client))

	if readOnly {
		return
	}

	s.AddTool(mcp.NewTool("add_reaction",
		mcp.WithDescription("Add an emoji reaction to a message"),
		mcp.WithString("channel_id", mcp.Required(), mcp.Description("Channel ID")),
		mcp.WithString("timestamp", mcp.Required(), mcp.Description("Message timestamp")),
		mcp.WithString("name", mcp.Required(), mcp.Description("Emoji name (without colons)")),
	), makeAddReaction(client))

	s.AddTool(mcp.NewTool("remove_reaction",
		mcp.WithDescription("Remove an emoji reaction from a message"),
		mcp.WithString("channel_id", mcp.Required(), mcp.Description("Channel ID")),
		mcp.WithString("timestamp", mcp.Required(), mcp.Description("Message timestamp")),
		mcp.WithString("name", mcp.Required(), mcp.Description("Emoji name (without colons)")),
	), makeRemoveReaction(client))
}

func makeAddReaction(client slack.Service) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		channelID, err := request.RequireString("channel_id")
		if err != nil {
			return errResult(err), nil
		}
		timestamp, err := request.RequireString("timestamp")
		if err != nil {
			return errResult(err), nil
		}
		name, err := request.RequireString("name")
		if err != nil {
			return errResult(err), nil
		}
		if err := client.AddReaction(channelID, timestamp, name); err != nil {
			return errResult(err), nil
		}
		return mcp.NewToolResultText(toJSON(map[string]string{"status": "added", "reaction": name})), nil
	}
}

func makeRemoveReaction(client slack.Service) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		channelID, err := request.RequireString("channel_id")
		if err != nil {
			return errResult(err), nil
		}
		timestamp, err := request.RequireString("timestamp")
		if err != nil {
			return errResult(err), nil
		}
		name, err := request.RequireString("name")
		if err != nil {
			return errResult(err), nil
		}
		if err := client.RemoveReaction(channelID, timestamp, name); err != nil {
			return errResult(err), nil
		}
		return mcp.NewToolResultText(toJSON(map[string]string{"status": "removed", "reaction": name})), nil
	}
}

func makeListReactions(client slack.Service) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		userID := request.GetString("user_id", "")
		limit := request.GetInt("limit", 100)

		result, err := client.ListReactions(userID, slack.PaginationParams{Limit: limit})
		if err != nil {
			return errResult(err), nil
		}
		return mcp.NewToolResultText(toJSON(result)), nil
	}
}
