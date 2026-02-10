package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/jackchuka/slackcli/internal/slack"
)

func registerSearchTools(s *server.MCPServer, client slack.Service) {
	s.AddTool(mcp.NewTool("search_messages",
		mcp.WithDescription("Search for messages in Slack"),
		mcp.WithString("query", mcp.Required(), mcp.Description("Search query")),
		mcp.WithString("sort", mcp.Description("Sort field: timestamp or score"), mcp.DefaultString("timestamp")),
		mcp.WithString("sort_dir", mcp.Description("Sort direction: asc or desc"), mcp.DefaultString("desc")),
		mcp.WithNumber("limit", mcp.Description("Max results to return"), mcp.DefaultNumber(20)),
	), makeSearchMessages(client))
}

func makeSearchMessages(client slack.Service) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		query, err := request.RequireString("query")
		if err != nil {
			return errResult(err), nil
		}
		sort := request.GetString("sort", "timestamp")
		sortDir := request.GetString("sort_dir", "desc")
		limit := request.GetInt("limit", 20)

		result, err := client.SearchMessages(slack.SearchParams{
			Query:      query,
			Sort:       sort,
			SortDir:    sortDir,
			Pagination: slack.PaginationParams{Limit: limit},
		})
		if err != nil {
			return errResult(err), nil
		}
		return mcp.NewToolResultText(toJSON(result)), nil
	}
}
