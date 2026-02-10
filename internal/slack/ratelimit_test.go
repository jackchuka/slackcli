package slack

import (
	"errors"
	"testing"
	"time"

	slackapi "github.com/slack-go/slack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRetry(t *testing.T) {
	t.Run("succeeds on first try", func(t *testing.T) {
		calls := 0
		result, err := retry(func() (string, error) {
			calls++
			return "ok", nil
		})

		require.NoError(t, err)
		assert.Equal(t, "ok", result)
		assert.Equal(t, 1, calls)
	})

	t.Run("succeeds after rate limit retries", func(t *testing.T) {
		calls := 0
		result, err := retry(func() (string, error) {
			calls++
			if calls < 3 {
				return "", &slackapi.RateLimitedError{RetryAfter: time.Millisecond}
			}
			return "recovered", nil
		})

		require.NoError(t, err)
		assert.Equal(t, "recovered", result)
		assert.Equal(t, 3, calls)
	})

	t.Run("fails after max retries", func(t *testing.T) {
		calls := 0
		result, err := retry(func() (string, error) {
			calls++
			return "", &slackapi.RateLimitedError{RetryAfter: time.Millisecond}
		})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "rate limited after")
		assert.Empty(t, result)
		assert.Equal(t, maxRetries+1, calls)
	})

	t.Run("non-rate-limit error fails immediately", func(t *testing.T) {
		calls := 0
		result, err := retry(func() (string, error) {
			calls++
			return "", errors.New("bad request")
		})

		require.Error(t, err)
		assert.Equal(t, "bad request", err.Error())
		assert.Empty(t, result)
		assert.Equal(t, 1, calls)
	})

	t.Run("works with struct return type", func(t *testing.T) {
		type data struct {
			Value int
		}
		result, err := retry(func() (data, error) {
			return data{Value: 42}, nil
		})

		require.NoError(t, err)
		assert.Equal(t, 42, result.Value)
	})

	t.Run("works with pointer return type", func(t *testing.T) {
		type data struct {
			Value int
		}
		result, err := retry(func() (*data, error) {
			return &data{Value: 99}, nil
		})

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 99, result.Value)
	})
}
