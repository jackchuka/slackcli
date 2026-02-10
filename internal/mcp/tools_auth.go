package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/jackchuka/slackcli/internal/slack"
)

func registerAuthTools(s *server.MCPServer, client slack.Service) {
	s.AddTool(mcp.NewTool("auth_test",
		mcp.WithDescription("Test authentication and get current user info"),
	), makeAuthTest(client))
}

func makeAuthTest(client slack.Service) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		result, err := client.AuthTest()
		if err != nil {
			return errResult(err), nil
		}
		return mcp.NewToolResultText(toJSON(result)), nil
	}
}
