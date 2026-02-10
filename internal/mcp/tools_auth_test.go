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

func TestMakeAuthTest(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mock := mocks.NewMockService(ctrl)

		mock.EXPECT().AuthTest().Return(&slack.AuthTestResult{
			UserID: "U123",
			User:   "alice",
			TeamID: "T456",
			Team:   "myteam",
			URL:    "https://myteam.slack.com/",
		}, nil)

		handler := makeAuthTest(mock)
		result, err := handler(context.Background(), newRequest(nil))

		require.NoError(t, err)
		assert.False(t, result.IsError)
	})

	t.Run("error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mock := mocks.NewMockService(ctrl)

		mock.EXPECT().AuthTest().Return(nil, &slack.SlackError{
			Code: slack.ErrAuth, Message: "invalid_auth",
		})

		handler := makeAuthTest(mock)
		result, err := handler(context.Background(), newRequest(nil))

		require.NoError(t, err)
		assert.True(t, result.IsError)
	})
}
