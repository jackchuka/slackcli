package main

import (
	"errors"
	"fmt"
	"testing"

	"github.com/jackchuka/slackcli/internal/slack"
	"github.com/stretchr/testify/assert"
)

func TestExitCode(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want int
	}{
		{
			name: "auth error returns 2",
			err:  &slack.SlackError{Code: slack.ErrAuth, Message: "invalid_auth"},
			want: 2,
		},
		{
			name: "not found returns 3",
			err:  &slack.SlackError{Code: slack.ErrNotFound, Message: "channel_not_found"},
			want: 3,
		},
		{
			name: "other SlackError returns 1",
			err:  &slack.SlackError{Code: slack.ErrAPI, Message: "unknown"},
			want: 1,
		},
		{
			name: "non-SlackError returns 1",
			err:  errors.New("random error"),
			want: 1,
		},
		{
			name: "wrapped auth error returns 2",
			err:  fmt.Errorf("wrap: %w", &slack.SlackError{Code: slack.ErrAuth, Message: "not_authed"}),
			want: 2,
		},
		{
			name: "wrapped not_found error returns 3",
			err:  fmt.Errorf("wrap: %w", &slack.SlackError{Code: slack.ErrNotFound, Message: "user_not_found"}),
			want: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, exitCode(tt.err))
		})
	}
}
