package auth

import "os"

// TokenStore defines the interface for storing and retrieving tokens.
type TokenStore interface {
	Get(workspace string) (string, error)
	Set(workspace, token string) error
	Delete(workspace string) error
}

// EnvProvider returns the token from environment variables.
type EnvProvider struct{}

func (e *EnvProvider) Token() string {
	return os.Getenv("SLACK_TOKEN")
}

// Resolver resolves a token through the chain: flag -> env -> config.
type Resolver struct {
	FlagToken string
	Env       *EnvProvider
	ConfigFn  func() string
}

func NewResolver(flagToken string, configFn func() string) *Resolver {
	return &Resolver{
		FlagToken: flagToken,
		Env:       &EnvProvider{},
		ConfigFn:  configFn,
	}
}

func (r *Resolver) Resolve() string {
	if r.FlagToken != "" {
		return r.FlagToken
	}
	if t := r.Env.Token(); t != "" {
		return t
	}
	if r.ConfigFn != nil {
		return r.ConfigFn()
	}
	return ""
}
