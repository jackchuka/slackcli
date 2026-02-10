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

func TestMakeListUsers(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mock := mocks.NewMockService(ctrl)

		mock.EXPECT().ListUsers(slack.PaginationParams{
			Limit: 100,
			All:   false,
		}).Return(&slack.PaginatedResult[slack.User]{
			Items: []slack.User{{ID: "U1", Name: "alice"}},
		}, nil)

		handler := makeListUsers(mock)
		result, err := handler(context.Background(), newRequest(nil))

		require.NoError(t, err)
		assert.False(t, result.IsError)
	})
}

func TestMakeGetUserInfo(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mock := mocks.NewMockService(ctrl)

		mock.EXPECT().GetUserInfo("U123").Return(&slack.User{
			ID: "U123", Name: "alice", RealName: "Alice Smith",
		}, nil)

		handler := makeGetUserInfo(mock)
		result, err := handler(context.Background(), newRequest(map[string]any{
			"user_id": "U123",
		}))

		require.NoError(t, err)
		assert.False(t, result.IsError)
	})

	t.Run("missing user_id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mock := mocks.NewMockService(ctrl)

		handler := makeGetUserInfo(mock)
		result, err := handler(context.Background(), newRequest(nil))

		require.NoError(t, err)
		assert.True(t, result.IsError)
	})
}

func TestMakeGetUserPresence(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mock := mocks.NewMockService(ctrl)

		mock.EXPECT().GetUserPresence("U123").Return("active", nil)

		handler := makeGetUserPresence(mock)
		result, err := handler(context.Background(), newRequest(map[string]any{
			"user_id": "U123",
		}))

		require.NoError(t, err)
		assert.False(t, result.IsError)
	})
}
