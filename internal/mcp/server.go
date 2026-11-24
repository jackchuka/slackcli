package mcp

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/jackchuka/slackcli/internal/slack"
	"github.com/jackchuka/slackcli/internal/version"
)

func NewServer(client slack.Service, readOnly bool) *server.MCPServer {
	s := server.NewMCPServer(
		"slackcli",
		version.Version,
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

func readOnlyHint() mcp.ToolOption {
	return func(t *mcp.Tool) {
		mcp.WithReadOnlyHintAnnotation(true)(t)
		mcp.WithDestructiveHintAnnotation(false)(t)
	}
}

// additiveHint marks a tool that writes without overwriting or removing anything.
func additiveHint() mcp.ToolOption {
	return mcp.WithDestructiveHintAnnotation(false)
}
