package slack

import (
	"testing"

	slackapi "github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
)

func TestChannelFromAPI(t *testing.T) {
	input := slackapi.Channel{
		GroupConversation: slackapi.GroupConversation{
			Conversation: slackapi.Conversation{
				ID:         "C123ABC",
				NumMembers: 42,
				IsPrivate:  true,
			},
			Name:       "general",
			IsArchived: true,
			Topic:      slackapi.Topic{Value: "team discussion"},
			Purpose:    slackapi.Purpose{Value: "general chat"},
		},
		IsMember: true,
	}

	got := channelFromAPI(input)

	assert.Equal(t, "C123ABC", got.ID)
	assert.Equal(t, "general", got.Name)
	assert.Equal(t, "team discussion", got.Topic)
	assert.Equal(t, "general chat", got.Purpose)
	assert.Equal(t, 42, got.NumMembers)
	assert.True(t, got.IsArchived)
	assert.True(t, got.IsPrivate)
	assert.True(t, got.IsMember)
}

func TestChannelFromAPI_Empty(t *testing.T) {
	got := channelFromAPI(slackapi.Channel{})

	assert.Empty(t, got.ID)
	assert.Empty(t, got.Name)
	assert.Empty(t, got.Topic)
	assert.Empty(t, got.Purpose)
	assert.Zero(t, got.NumMembers)
	assert.False(t, got.IsArchived)
	assert.False(t, got.IsPrivate)
	assert.False(t, got.IsMember)
	assert.Zero(t, got.Created)
}
