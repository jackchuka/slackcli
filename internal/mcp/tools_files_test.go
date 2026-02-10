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

func TestMakeListFiles(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mock := mocks.NewMockService(ctrl)

		mock.EXPECT().ListFiles(
			slack.PaginationParams{Limit: 100}, "", "",
		).Return(&slack.PaginatedResult[slack.File]{
			Items: []slack.File{{ID: "F1", Name: "doc.pdf"}},
		}, nil)

		handler := makeListFiles(mock)
		result, err := handler(context.Background(), newRequest(nil))

		require.NoError(t, err)
		assert.False(t, result.IsError)
	})

	t.Run("with filters", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mock := mocks.NewMockService(ctrl)

		mock.EXPECT().ListFiles(
			slack.PaginationParams{Limit: 50}, "C123", "U456",
		).Return(&slack.PaginatedResult[slack.File]{Items: nil}, nil)

		handler := makeListFiles(mock)
		result, err := handler(context.Background(), newRequest(map[string]any{
			"channel_id": "C123",
			"user_id":    "U456",
			"limit":      float64(50),
		}))

		require.NoError(t, err)
		assert.False(t, result.IsError)
	})
}

func TestMakeGetFileInfo(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mock := mocks.NewMockService(ctrl)

		mock.EXPECT().GetFileInfo("F123").Return(&slack.File{
			ID: "F123", Name: "report.pdf",
		}, nil)

		handler := makeGetFileInfo(mock)
		result, err := handler(context.Background(), newRequest(map[string]any{
			"file_id": "F123",
		}))

		require.NoError(t, err)
		assert.False(t, result.IsError)
	})

	t.Run("missing file_id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mock := mocks.NewMockService(ctrl)

		handler := makeGetFileInfo(mock)
		result, err := handler(context.Background(), newRequest(nil))

		require.NoError(t, err)
		assert.True(t, result.IsError)
	})
}

func TestMakeDeleteFile(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mock := mocks.NewMockService(ctrl)

		mock.EXPECT().DeleteFile("F123").Return(nil)

		handler := makeDeleteFile(mock)
		result, err := handler(context.Background(), newRequest(map[string]any{
			"file_id": "F123",
		}))

		require.NoError(t, err)
		assert.False(t, result.IsError)
	})
}
