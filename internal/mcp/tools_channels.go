package mcp

import (
	"context"
	"encoding/json"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/jackchuka/slackcli/internal/slack"
)

func registerChannelTools(s *server.MCPServer, client slack.Service, readOnly bool) {
	s.AddTool(mcp.NewTool("list_channels",
		mcp.WithDescription("List Slack channels"),
		mcp.WithNumber("limit", mcp.Description("Max channels to return"), mcp.DefaultNumber(100)),
		mcp.WithBoolean("all", mcp.Description("Fetch all channels (auto-paginate)")),
		mcp.WithString("cursor", mcp.Description("Pagination cursor")),
	), makeListChannels(client))

	s.AddTool(mcp.NewTool("get_channel_info",
		mcp.WithDescription("Get information about a Slack channel"),
		mcp.WithString("channel_id", mcp.Required(), mcp.Description("Channel ID")),
	), makeGetChannelInfo(client))

	if readOnly {
		return
	}

	s.AddTool(mcp.NewTool("create_channel",
		mcp.WithDescription("Create a new Slack channel"),
		mcp.WithString("name", mcp.Required(), mcp.Description("Channel name")),
		mcp.WithBoolean("is_private", mcp.Description("Create as private channel")),
	), makeCreateChannel(client))

	s.AddTool(mcp.NewTool("archive_channel",
		mcp.WithDescription("Archive a Slack channel"),
		mcp.WithString("channel_id", mcp.Required(), mcp.Description("Channel ID")),
	), makeArchiveChannel(client))

	s.AddTool(mcp.NewTool("invite_to_channel",
		mcp.WithDescription("Invite users to a channel"),
		mcp.WithString("channel_id", mcp.Required(), mcp.Description("Channel ID")),
		mcp.WithArray("user_ids", mcp.Required(), mcp.Description("User IDs to invite")),
	), makeInviteToChannel(client))

	s.AddTool(mcp.NewTool("kick_from_channel",
		mcp.WithDescription("Remove a user from a channel"),
		mcp.WithString("channel_id", mcp.Required(), mcp.Description("Channel ID")),
		mcp.WithString("user_id", mcp.Required(), mcp.Description("User ID to remove")),
	), makeKickFromChannel(client))

	s.AddTool(mcp.NewTool("set_channel_topic",
		mcp.WithDescription("Set a channel's topic"),
		mcp.WithString("channel_id", mcp.Required(), mcp.Description("Channel ID")),
		mcp.WithString("topic", mcp.Required(), mcp.Description("New topic")),
	), makeSetChannelTopic(client))

	s.AddTool(mcp.NewTool("set_channel_purpose",
		mcp.WithDescription("Set a channel's purpose"),
		mcp.WithString("channel_id", mcp.Required(), mcp.Description("Channel ID")),
		mcp.WithString("purpose", mcp.Required(), mcp.Description("New purpose")),
	), makeSetChannelPurpose(client))
}

func toJSON(v any) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}

func makeListChannels(client slack.Service) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		limit := request.GetInt("limit", 100)
		all := request.GetBool("all", false)
		cursor := request.GetString("cursor", "")

		result, err := client.ListChannels(slack.PaginationParams{
			Cursor: cursor,
			Limit:  limit,
			All:    all,
		})
		if err != nil {
			return errResult(err), nil
		}
		return mcp.NewToolResultText(toJSON(result)), nil
	}
}

func makeGetChannelInfo(client slack.Service) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		channelID, err := request.RequireString("channel_id")
		if err != nil {
			return errResult(err), nil
		}
		ch, err := client.GetChannelInfo(channelID)
		if err != nil {
			return errResult(err), nil
		}
		return mcp.NewToolResultText(toJSON(ch)), nil
	}
}

func makeCreateChannel(client slack.Service) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name, err := request.RequireString("name")
		if err != nil {
			return errResult(err), nil
		}
		isPrivate := request.GetBool("is_private", false)
		ch, err := client.CreateChannel(name, isPrivate)
		if err != nil {
			return errResult(err), nil
		}
		return mcp.NewToolResultText(toJSON(ch)), nil
	}
}

func makeArchiveChannel(client slack.Service) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		channelID, err := request.RequireString("channel_id")
		if err != nil {
			return errResult(err), nil
		}
		if err := client.ArchiveChannel(channelID); err != nil {
			return errResult(err), nil
		}
		return mcp.NewToolResultText(toJSON(map[string]string{"status": "archived", "channel_id": channelID})), nil
	}
}

func makeInviteToChannel(client slack.Service) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		channelID, err := request.RequireString("channel_id")
		if err != nil {
			return errResult(err), nil
		}
		userIDs := request.GetStringSlice("user_ids", nil)
		if len(userIDs) == 0 {
			return mcp.NewToolResultError("user_ids is required"), nil
		}
		if err := client.InviteToChannel(channelID, userIDs...); err != nil {
			return errResult(err), nil
		}
		return mcp.NewToolResultText(toJSON(map[string]any{"status": "invited", "channel_id": channelID, "user_ids": userIDs})), nil
	}
}

func makeKickFromChannel(client slack.Service) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		channelID, err := request.RequireString("channel_id")
		if err != nil {
			return errResult(err), nil
		}
		userID, err := request.RequireString("user_id")
		if err != nil {
			return errResult(err), nil
		}
		if err := client.KickFromChannel(channelID, userID); err != nil {
			return errResult(err), nil
		}
		return mcp.NewToolResultText(toJSON(map[string]string{"status": "removed", "channel_id": channelID, "user_id": userID})), nil
	}
}

func makeSetChannelTopic(client slack.Service) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		channelID, err := request.RequireString("channel_id")
		if err != nil {
			return errResult(err), nil
		}
		topic, err := request.RequireString("topic")
		if err != nil {
			return errResult(err), nil
		}
		if err := client.SetChannelTopic(channelID, topic); err != nil {
			return errResult(err), nil
		}
		return mcp.NewToolResultText(toJSON(map[string]string{"status": "updated", "channel_id": channelID, "topic": topic})), nil
	}
}

func makeSetChannelPurpose(client slack.Service) server.ToolHandlerFunc {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		channelID, err := request.RequireString("channel_id")
		if err != nil {
			return errResult(err), nil
		}
		purpose, err := request.RequireString("purpose")
		if err != nil {
			return errResult(err), nil
		}
		if err := client.SetChannelPurpose(channelID, purpose); err != nil {
			return errResult(err), nil
		}
		return mcp.NewToolResultText(toJSON(map[string]string{"status": "updated", "channel_id": channelID, "purpose": purpose})), nil
	}
}
