package slack

import (
	"testing"

	slackapi "github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
)

func TestUserFromAPI(t *testing.T) {
	input := slackapi.User{
		ID:       "U123ABC",
		Name:     "jdoe",
		RealName: "Jane Doe",
		Profile: slackapi.UserProfile{
			Email: "jane@example.com",
		},
		IsAdmin: true,
		IsBot:   false,
		Deleted: false,
		TZ:      "America/New_York",
	}

	got := userFromAPI(input)

	assert.Equal(t, "U123ABC", got.ID)
	assert.Equal(t, "jdoe", got.Name)
	assert.Equal(t, "Jane Doe", got.RealName)
	assert.Equal(t, "jane@example.com", got.Email)
	assert.True(t, got.IsAdmin)
	assert.False(t, got.IsBot)
	assert.False(t, got.Deleted)
	assert.Equal(t, "America/New_York", got.TZ)
	// Presence is not set by userFromAPI
	assert.Empty(t, got.Presence)
}

func TestUserFromAPI_Empty(t *testing.T) {
	got := userFromAPI(slackapi.User{})

	assert.Empty(t, got.ID)
	assert.Empty(t, got.Name)
	assert.Empty(t, got.RealName)
	assert.Empty(t, got.Email)
	assert.False(t, got.IsAdmin)
	assert.False(t, got.IsBot)
	assert.False(t, got.Deleted)
	assert.Empty(t, got.TZ)
}
