package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvProvider_Token(t *testing.T) {
	t.Run("returns SLACK_TOKEN when set", func(t *testing.T) {
		t.Setenv("SLACK_TOKEN", "xoxb-test-token")
		p := &EnvProvider{}
		assert.Equal(t, "xoxb-test-token", p.Token())
	})

	t.Run("returns empty when SLACK_TOKEN not set", func(t *testing.T) {
		t.Setenv("SLACK_TOKEN", "")
		p := &EnvProvider{}
		assert.Empty(t, p.Token())
	})
}

func TestResolver_Resolve(t *testing.T) {
	t.Run("flag token takes priority", func(t *testing.T) {
		t.Setenv("SLACK_TOKEN", "env-token")
		r := NewResolver("flag-token", func() string { return "config-token" })
		assert.Equal(t, "flag-token", r.Resolve())
	})

	t.Run("env token when no flag", func(t *testing.T) {
		t.Setenv("SLACK_TOKEN", "env-token")
		r := NewResolver("", func() string { return "config-token" })
		assert.Equal(t, "env-token", r.Resolve())
	})

	t.Run("config token when no flag or env", func(t *testing.T) {
		t.Setenv("SLACK_TOKEN", "")
		r := NewResolver("", func() string { return "config-token" })
		assert.Equal(t, "config-token", r.Resolve())
	})

	t.Run("returns empty when all sources empty", func(t *testing.T) {
		t.Setenv("SLACK_TOKEN", "")
		r := NewResolver("", func() string { return "" })
		assert.Empty(t, r.Resolve())
	})

	t.Run("nil ConfigFn returns empty", func(t *testing.T) {
		t.Setenv("SLACK_TOKEN", "")
		r := &Resolver{Env: &EnvProvider{}}
		assert.Empty(t, r.Resolve())
	})
}
