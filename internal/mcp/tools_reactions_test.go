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

func TestMakeAddReaction(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mock := mocks.NewMockService(ctrl)

		mock.EXPECT().AddReaction("C123", "1234.5678", "thumbsup").Return(nil)

		handler := makeAddReaction(mock)
		result, err := handler(context.Background(), newRequest(map[string]any{
			"channel_id": "C123",
			"timestamp":  "1234.5678",
			"name":       "thumbsup",
		}))

		require.NoError(t, err)
		assert.False(t, result.IsError)
	})

	t.Run("error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mock := mocks.NewMockService(ctrl)

		mock.EXPECT().AddReaction("C123", "1234.5678", "thumbsup").Return(
			&slack.SlackError{Code: slack.ErrNotFound, Message: "message_not_found"},
		)

		handler := makeAddReaction(mock)
		result, err := handler(context.Background(), newRequest(map[string]any{
			"channel_id": "C123",
			"timestamp":  "1234.5678",
			"name":       "thumbsup",
		}))

		require.NoError(t, err)
		assert.True(t, result.IsError)
	})
}

func TestMakeRemoveReaction(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mock := mocks.NewMockService(ctrl)

		mock.EXPECT().RemoveReaction("C123", "1234.5678", "thumbsup").Return(nil)

		handler := makeRemoveReaction(mock)
		result, err := handler(context.Background(), newRequest(map[string]any{
			"channel_id": "C123",
			"timestamp":  "1234.5678",
			"name":       "thumbsup",
		}))

		require.NoError(t, err)
		assert.False(t, result.IsError)
	})
}

func TestMakeListReactions(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mock := mocks.NewMockService(ctrl)

		mock.EXPECT().ListReactions("U123", slack.PaginationParams{Limit: 100}).Return(
			&slack.PaginatedResult[slack.ReactedItem]{
				Items: []slack.ReactedItem{{Type: "message"}},
			}, nil)

		handler := makeListReactions(mock)
		result, err := handler(context.Background(), newRequest(map[string]any{
			"user_id": "U123",
		}))

		require.NoError(t, err)
		assert.False(t, result.IsError)
	})
}
