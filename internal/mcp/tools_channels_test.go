package mcp

import (
	"context"
	"testing"

	"github.com/jackchuka/slackcli/internal/slack"
	"github.com/jackchuka/slackcli/internal/slack/mocks"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func newRequest(args map[string]any) mcp.CallToolRequest {
	return mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Arguments: args,
		},
	}
}

func TestMakeListChannels(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mock := mocks.NewMockService(ctrl)

		mock.EXPECT().ListChannels(slack.PaginationParams{
			Cursor: "",
			Limit:  100,
			All:    false,
		}).Return(&slack.PaginatedResult[slack.Channel]{
			Items: []slack.Channel{{ID: "C1", Name: "general"}},
		}, nil)

		handler := makeListChannels(mock)
		result, err := handler(context.Background(), newRequest(nil))

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.False(t, result.IsError)
	})

	t.Run("error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mock := mocks.NewMockService(ctrl)

		mock.EXPECT().ListChannels(gomock.Any()).Return(nil, &slack.SlackError{
			Code: slack.ErrAuth, Message: "invalid_auth",
		})

		handler := makeListChannels(mock)
		result, err := handler(context.Background(), newRequest(nil))

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.True(t, result.IsError)
	})
}

func TestMakeGetChannelInfo(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mock := mocks.NewMockService(ctrl)

		mock.EXPECT().GetChannelInfo("C123").Return(&slack.Channel{
			ID: "C123", Name: "general",
		}, nil)

		handler := makeGetChannelInfo(mock)
		result, err := handler(context.Background(), newRequest(map[string]any{
			"channel_id": "C123",
		}))

		require.NoError(t, err)
		assert.False(t, result.IsError)
	})

	t.Run("missing channel_id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mock := mocks.NewMockService(ctrl)

		handler := makeGetChannelInfo(mock)
		result, err := handler(context.Background(), newRequest(nil))

		require.NoError(t, err)
		assert.True(t, result.IsError)
	})
}

func TestMakeCreateChannel(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mock := mocks.NewMockService(ctrl)

		mock.EXPECT().CreateChannel("dev", false).Return(&slack.Channel{
			ID: "C456", Name: "dev",
		}, nil)

		handler := makeCreateChannel(mock)
		result, err := handler(context.Background(), newRequest(map[string]any{
			"name": "dev",
		}))

		require.NoError(t, err)
		assert.False(t, result.IsError)
	})
}

func TestMakeArchiveChannel(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mock := mocks.NewMockService(ctrl)

		mock.EXPECT().ArchiveChannel("C123").Return(nil)

		handler := makeArchiveChannel(mock)
		result, err := handler(context.Background(), newRequest(map[string]any{
			"channel_id": "C123",
		}))

		require.NoError(t, err)
		assert.False(t, result.IsError)
	})
}

func TestMakeSetChannelTopic(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mock := mocks.NewMockService(ctrl)

		mock.EXPECT().SetChannelTopic("C123", "new topic").Return(nil)

		handler := makeSetChannelTopic(mock)
		result, err := handler(context.Background(), newRequest(map[string]any{
			"channel_id": "C123",
			"topic":      "new topic",
		}))

		require.NoError(t, err)
		assert.False(t, result.IsError)
	})
}

func TestMakeSetChannelPurpose(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mock := mocks.NewMockService(ctrl)

		mock.EXPECT().SetChannelPurpose("C123", "new purpose").Return(nil)

		handler := makeSetChannelPurpose(mock)
		result, err := handler(context.Background(), newRequest(map[string]any{
			"channel_id": "C123",
			"purpose":    "new purpose",
		}))

		require.NoError(t, err)
		assert.False(t, result.IsError)
	})
}
