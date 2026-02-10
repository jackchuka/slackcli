package slack

import (
	"testing"

	slackapi "github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConvertReactedItems(t *testing.T) {
	t.Run("nil input returns nil", func(t *testing.T) {
		assert.Nil(t, convertReactedItems(nil))
	})

	t.Run("empty slice returns nil", func(t *testing.T) {
		assert.Nil(t, convertReactedItems([]slackapi.ReactedItem{}))
	})

	t.Run("message with reactions", func(t *testing.T) {
		input := []slackapi.ReactedItem{
			{
				Item: slackapi.Item{
					Type: "message",
					Message: &slackapi.Message{
						Msg: slackapi.Msg{
							Channel:   "C123",
							Timestamp: "1700000000.000100",
						},
					},
				},
				Reactions: []slackapi.ItemReaction{
					{Name: "thumbsup", Count: 3, Users: []string{"U1", "U2", "U3"}},
					{Name: "heart", Count: 1, Users: []string{"U4"}},
				},
			},
		}

		got := convertReactedItems(input)

		require.Len(t, got, 1)
		assert.Equal(t, "message", got[0].Type)
		assert.Equal(t, "C123", got[0].Channel)
		assert.Equal(t, "1700000000.000100", got[0].Timestamp)
		require.Len(t, got[0].Reactions, 2)
		assert.Equal(t, "thumbsup", got[0].Reactions[0].Name)
		assert.Equal(t, 3, got[0].Reactions[0].Count)
		assert.Equal(t, []string{"U1", "U2", "U3"}, got[0].Reactions[0].Users)
		assert.Equal(t, "heart", got[0].Reactions[1].Name)
	})

	t.Run("nil Message leaves Channel and Timestamp empty", func(t *testing.T) {
		input := []slackapi.ReactedItem{
			{
				Item: slackapi.Item{
					Type:    "file",
					Message: nil,
				},
				Reactions: []slackapi.ItemReaction{
					{Name: "eyes", Count: 1, Users: []string{"U1"}},
				},
			},
		}

		got := convertReactedItems(input)

		require.Len(t, got, 1)
		assert.Equal(t, "file", got[0].Type)
		assert.Empty(t, got[0].Channel)
		assert.Empty(t, got[0].Timestamp)
		require.Len(t, got[0].Reactions, 1)
	})

	t.Run("item with no reactions", func(t *testing.T) {
		input := []slackapi.ReactedItem{
			{
				Item: slackapi.Item{
					Type: "message",
					Message: &slackapi.Message{
						Msg: slackapi.Msg{Channel: "C1", Timestamp: "1.0"},
					},
				},
				Reactions: nil,
			},
		}

		got := convertReactedItems(input)

		require.Len(t, got, 1)
		assert.Nil(t, got[0].Reactions)
	})
}
