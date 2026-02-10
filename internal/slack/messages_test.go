package slack

import (
	"testing"
	"time"

	slackapi "github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
)

func TestMessageFromAPI(t *testing.T) {
	input := slackapi.Message{
		Msg: slackapi.Msg{
			Timestamp:       "1700000000.000100",
			User:            "U123ABC",
			Text:            "hello world",
			ThreadTimestamp: "1700000000.000001",
			Type:            "message",
		},
	}

	got := messageFromAPI(input)

	assert.Equal(t, "1700000000.000100", got.Timestamp)
	assert.Equal(t, "U123ABC", got.User)
	assert.Equal(t, "hello world", got.Text)
	assert.Equal(t, "1700000000.000001", got.ThreadTS)
	assert.Equal(t, "message", got.Type)
	// Channel is NOT set by messageFromAPI — set separately in listMessagesPage
	assert.Empty(t, got.Channel)
}

func TestMessageFromAPI_Empty(t *testing.T) {
	got := messageFromAPI(slackapi.Message{})

	assert.Empty(t, got.Timestamp)
	assert.Empty(t, got.User)
	assert.Empty(t, got.Text)
	assert.Empty(t, got.ThreadTS)
	assert.Empty(t, got.Type)
}

func TestFormatTimestamp(t *testing.T) {
	tests := []struct {
		name string
		time time.Time
		want string
	}{
		{"unix epoch", time.Unix(0, 0), "0.000000"},
		{"specific time", time.Unix(1700000000, 0), "1700000000.000000"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, formatTimestamp(tt.time))
		})
	}
}
