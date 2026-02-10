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

func TestMakeSearchMessages(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mock := mocks.NewMockService(ctrl)

		mock.EXPECT().SearchMessages(slack.SearchParams{
			Query:      "important",
			Sort:       "timestamp",
			SortDir:    "desc",
			Pagination: slack.PaginationParams{Limit: 20},
		}).Return(&slack.SearchResult{
			Matches: []slack.Message{{Text: "this is important"}},
			Total:   1,
		}, nil)

		handler := makeSearchMessages(mock)
		result, err := handler(context.Background(), newRequest(map[string]any{
			"query": "important",
		}))

		require.NoError(t, err)
		assert.False(t, result.IsError)
	})

	t.Run("missing query", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mock := mocks.NewMockService(ctrl)

		handler := makeSearchMessages(mock)
		result, err := handler(context.Background(), newRequest(nil))

		require.NoError(t, err)
		assert.True(t, result.IsError)
	})
}
