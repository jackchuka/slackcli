package slack

import (
	"fmt"
	"time"

	slackapi "github.com/slack-go/slack"
)

const maxRetries = 3

func retry[T any](fn func() (T, error)) (T, error) {
	var zero T
	for attempt := 0; attempt <= maxRetries; attempt++ {
		result, err := fn()
		if err == nil {
			return result, nil
		}
		rateLimitErr, ok := err.(*slackapi.RateLimitedError)
		if !ok {
			return zero, err
		}
		if attempt == maxRetries {
			return zero, fmt.Errorf("rate limited after %d retries: %w", maxRetries, err)
		}
		wait := rateLimitErr.RetryAfter
		if wait == 0 {
			wait = time.Duration(attempt+1) * time.Second
		}
		time.Sleep(wait)
	}
	return zero, fmt.Errorf("max retries exceeded")
}
