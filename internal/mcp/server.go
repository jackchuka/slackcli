package mcp

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/jackchuka/slackcli/internal/slack"
)

func NewServer(client slack.Service, readOnly bool) *server.MCPServer {
	s := server.NewMCPServer(
		"slackcli",
		"1.0.0",
		server.WithToolCapabilities(true),
		server.WithRecovery(),
	)

	registerChannelTools(s, client, readOnly)
	registerMessageTools(s, client, readOnly)
	registerUserTools(s, client)
	registerReactionTools(s, client, readOnly)
	registerFileTools(s, client, readOnly)
	registerSearchTools(s, client)
	registerAuthTools(s, client)

	return s
}

func errResult(err error) *mcp.CallToolResult {
	return mcp.NewToolResultError(err.Error())
}
