package slack

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClassifyError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantCode ErrorCode
		wantMsg  string
	}{
		// auth errors
		{"invalid_auth", errors.New("invalid_auth"), ErrAuth, "invalid_auth"},
		{"not_authed", errors.New("not_authed"), ErrAuth, "not_authed"},
		{"token_revoked", errors.New("token_revoked"), ErrAuth, "token_revoked"},
		{"token_expired", errors.New("token_expired"), ErrAuth, "token_expired"},
		{"account_inactive", errors.New("account_inactive"), ErrAuth, "account_inactive"},
		// not found errors
		{"channel_not_found", errors.New("channel_not_found"), ErrNotFound, "channel_not_found"},
		{"user_not_found", errors.New("user_not_found"), ErrNotFound, "user_not_found"},
		{"file_not_found", errors.New("file_not_found"), ErrNotFound, "file_not_found"},
		{"message_not_found", errors.New("message_not_found"), ErrNotFound, "message_not_found"},
		// permission errors
		{"not_in_channel", errors.New("not_in_channel"), ErrPermission, "not_in_channel"},
		{"missing_scope", errors.New("missing_scope"), ErrPermission, "missing_scope"},
		{"cannot_dm_bot", errors.New("cannot_dm_bot"), ErrPermission, "cannot_dm_bot"},
		{"restricted_action", errors.New("restricted_action"), ErrPermission, "restricted_action"},
		// validation errors
		{"too_many_attachments", errors.New("too_many_attachments"), ErrValidation, "too_many_attachments"},
		{"msg_too_long", errors.New("msg_too_long"), ErrValidation, "msg_too_long"},
		{"no_text", errors.New("no_text"), ErrValidation, "no_text"},
		{"invalid_blocks", errors.New("invalid_blocks"), ErrValidation, "invalid_blocks"},
		// default -> API error
		{"unknown error", errors.New("something_unexpected"), ErrAPI, "something_unexpected"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := classifyError(tt.err)
			require.NotNil(t, got)
			assert.Equal(t, tt.wantCode, got.Code)
			assert.Equal(t, tt.wantMsg, got.Message)
			assert.ErrorIs(t, got, tt.err)
		})
	}
}

func TestClassifyError_Nil(t *testing.T) {
	assert.Nil(t, classifyError(nil))
}

func TestSlackError_Error(t *testing.T) {
	tests := []struct {
		name   string
		err    SlackError
		expect string
	}{
		{
			name:   "without detail",
			err:    SlackError{Code: ErrAuth, Message: "invalid_auth"},
			expect: "auth_error: invalid_auth",
		},
		{
			name:   "with detail",
			err:    SlackError{Code: ErrNotFound, Message: "channel_not_found", Detail: "C123"},
			expect: "not_found: channel_not_found (C123)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expect, tt.err.Error())
		})
	}
}

func TestSlackError_Unwrap(t *testing.T) {
	original := errors.New("original error")
	se := &SlackError{Code: ErrAPI, Message: "test", Err: original}
	assert.Equal(t, original, se.Unwrap())

	seNilErr := &SlackError{Code: ErrAPI, Message: "test"}
	assert.Nil(t, seNilErr.Unwrap())
}
