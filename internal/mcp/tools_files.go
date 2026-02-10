package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/jackchuka/slackcli/internal/slack"
)

func registerFileTools(s *server.MCPServer, client slack.Service, readOnly bool) {
	s.AddTool(mcp.NewTool("list_files",
		mcp.WithDescription("List files in Slack"),
		mcp.WithString("channel_id", mcp.Description("Filter by channel ID")),
		mcp.WithString("user_id", mcp.Description("Filter by user ID")),
		mcp.WithNumber("limit", mcp.Description("Max files to return"), mcp.DefaultNumber(100)),
	), makeListFiles(client))

	s.AddTool(mcp.NewTool("get_file_info",
		mcp.WithDescription("Get information about a file"),
		mcp.WithString("file_id", mcp.Required(), mcp.Description("File ID")),
	), makeGetFileInfo(client))

	if readOnly {
		return
	}

	s.AddTool(mcp.NewTool("delete_file",
		mcp.WithDescription("Delete a file"),
		mcp.WithString("file_id", mcp.Required(), mcp.Description("File ID")),
	), makeDeleteFile(client))
}

func makeListFiles(client slack.Service) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		channelID := request.GetString("channel_id", "")
		userID := request.GetString("user_id", "")
		limit := request.GetInt("limit", 100)

		result, err := client.ListFiles(slack.PaginationParams{Limit: limit}, channelID, userID)
		if err != nil {
			return errResult(err), nil
		}
		return mcp.NewToolResultText(toJSON(result)), nil
	}
}

func makeGetFileInfo(client slack.Service) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		fileID, err := request.RequireString("file_id")
		if err != nil {
			return errResult(err), nil
		}
		file, err := client.GetFileInfo(fileID)
		if err != nil {
			return errResult(err), nil
		}
		return mcp.NewToolResultText(toJSON(file)), nil
	}
}

func makeDeleteFile(client slack.Service) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		fileID, err := request.RequireString("file_id")
		if err != nil {
			return errResult(err), nil
		}
		if err := client.DeleteFile(fileID); err != nil {
			return errResult(err), nil
		}
		return mcp.NewToolResultText(toJSON(map[string]string{"status": "deleted", "file_id": fileID})), nil
	}
}
