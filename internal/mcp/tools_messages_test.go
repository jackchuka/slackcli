package mcp

import (
	"context"
	"testing"

	"github.com/jackchuka/slackcli/internal/slack"
	"github.com/jackchuka/slackcli/internal/slack/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestMakeListMessages(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mock := mocks.NewMockService(ctrl)

		mock.EXPECT().ListMessages(slack.ListMessagesParams{
			ChannelID:  "C123",
			Pagination: slack.PaginationParams{Cursor: "", Limit: 100, All: false},
		}).Return(&slack.PaginatedResult[slack.Message]{
			Items: []slack.Message{{Text: "hello", User: "U1"}},
		}, nil)

		handler := makeListMessages(mock)
		result, err := handler(context.Background(), newRequest(map[string]any{
			"channel_id": "C123",
		}))

		require.NoError(t, err)
		assert.False(t, result.IsError)
	})

	t.Run("missing channel_id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mock := mocks.NewMockService(ctrl)

		handler := makeListMessages(mock)
		result, err := handler(context.Background(), newRequest(nil))

		require.NoError(t, err)
		assert.True(t, result.IsError)
	})
}

func TestMakeSendMessage(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mock := mocks.NewMockService(ctrl)

		mock.EXPECT().SendMessage(slack.SendMessageParams{
			ChannelID: "C123",
			Text:      "hello",
			ThreadTS:  "",
		}).Return(&slack.Message{
			Channel:   "C123",
			Text:      "hello",
			Timestamp: "1234.5678",
		}, nil)

		handler := makeSendMessage(mock)
		result, err := handler(context.Background(), newRequest(map[string]any{
			"channel_id": "C123",
			"text":       "hello",
		}))

		require.NoError(t, err)
		assert.False(t, result.IsError)
	})
}

func TestMakeEditMessage(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mock := mocks.NewMockService(ctrl)

		mock.EXPECT().EditMessage("C123", "1234.5678", "updated text").Return(&slack.Message{
			Channel: "C123", Timestamp: "1234.5678", Text: "updated text",
		}, nil)

		handler := makeEditMessage(mock)
		result, err := handler(context.Background(), newRequest(map[string]any{
			"channel_id": "C123",
			"timestamp":  "1234.5678",
			"text":       "updated text",
		}))

		require.NoError(t, err)
		assert.False(t, result.IsError)
	})
}

func TestMakeDeleteMessage(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mock := mocks.NewMockService(ctrl)

		mock.EXPECT().DeleteMessage("C123", "1234.5678").Return(nil)

		handler := makeDeleteMessage(mock)
		result, err := handler(context.Background(), newRequest(map[string]any{
			"channel_id": "C123",
			"timestamp":  "1234.5678",
		}))

		require.NoError(t, err)
		assert.False(t, result.IsError)
	})
}
