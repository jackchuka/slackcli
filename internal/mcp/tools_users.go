package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/jackchuka/slackcli/internal/slack"
)

func registerUserTools(s *server.MCPServer, client slack.Service) {
	s.AddTool(mcp.NewTool("list_users",
		mcp.WithDescription("List Slack users"),
		mcp.WithNumber("limit", mcp.Description("Max users to return"), mcp.DefaultNumber(100)),
		mcp.WithBoolean("all", mcp.Description("Fetch all users")),
	), makeListUsers(client))

	s.AddTool(mcp.NewTool("get_user_info",
		mcp.WithDescription("Get information about a Slack user"),
		mcp.WithString("user_id", mcp.Required(), mcp.Description("User ID")),
	), makeGetUserInfo(client))

	s.AddTool(mcp.NewTool("get_user_presence",
		mcp.WithDescription("Get a user's presence status"),
		mcp.WithString("user_id", mcp.Required(), mcp.Description("User ID")),
	), makeGetUserPresence(client))
}

func makeListUsers(client slack.Service) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		limit := request.GetInt("limit", 100)
		all := request.GetBool("all", false)

		result, err := client.ListUsers(slack.PaginationParams{Limit: limit, All: all})
		if err != nil {
			return errResult(err), nil
		}
		return mcp.NewToolResultText(toJSON(result)), nil
	}
}

func makeGetUserInfo(client slack.Service) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		userID, err := request.RequireString("user_id")
		if err != nil {
			return errResult(err), nil
		}
		user, err := client.GetUserInfo(userID)
		if err != nil {
			return errResult(err), nil
		}
		return mcp.NewToolResultText(toJSON(user)), nil
	}
}

func makeGetUserPresence(client slack.Service) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		userID, err := request.RequireString("user_id")
		if err != nil {
			return errResult(err), nil
		}
		presence, err := client.GetUserPresence(userID)
		if err != nil {
			return errResult(err), nil
		}
		return mcp.NewToolResultText(toJSON(map[string]string{"user_id": userID, "presence": presence})), nil
	}
}
