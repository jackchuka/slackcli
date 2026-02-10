package slack

import "fmt"

type ErrorCode string

const (
	ErrAuth       ErrorCode = "auth_error"
	ErrRateLimit  ErrorCode = "rate_limited"
	ErrNotFound   ErrorCode = "not_found"
	ErrPermission ErrorCode = "permission_denied"
	ErrValidation ErrorCode = "validation_error"
	ErrAPI        ErrorCode = "api_error"
	ErrNetwork    ErrorCode = "network_error"
)

type SlackError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Detail  string    `json:"detail,omitempty"`
	Err     error     `json:"-"`
}

func (e *SlackError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.Detail)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *SlackError) Unwrap() error {
	return e.Err
}

func classifyError(err error) *SlackError {
	if err == nil {
		return nil
	}
	msg := err.Error()
	switch msg {
	case "invalid_auth", "not_authed", "token_revoked", "token_expired", "account_inactive":
		return &SlackError{Code: ErrAuth, Message: msg, Err: err}
	case "channel_not_found", "user_not_found", "file_not_found", "message_not_found":
		return &SlackError{Code: ErrNotFound, Message: msg, Err: err}
	case "not_in_channel", "missing_scope", "cannot_dm_bot", "restricted_action":
		return &SlackError{Code: ErrPermission, Message: msg, Err: err}
	case "too_many_attachments", "msg_too_long", "no_text", "invalid_blocks":
		return &SlackError{Code: ErrValidation, Message: msg, Err: err}
	default:
		return &SlackError{Code: ErrAPI, Message: msg, Err: err}
	}
}
