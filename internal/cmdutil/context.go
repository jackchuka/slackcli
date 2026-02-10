package cmdutil

import (
	"context"

	"github.com/jackchuka/slackcli/internal/auth"
	"github.com/jackchuka/slackcli/internal/config"
	"github.com/jackchuka/slackcli/internal/output"
	"github.com/jackchuka/slackcli/internal/slack"
)

type contextKey string

const runContextKey contextKey = "run_context"

type RunContext struct {
	Config    *config.Config
	Client    slack.Service
	Formatter output.Formatter
	Writers   *output.Writers
	Resolver  *auth.Resolver
	ReadOnly  bool
}

func GetRunContext(ctx context.Context) *RunContext {
	return ctx.Value(runContextKey).(*RunContext)
}

func SetRunContext(ctx context.Context, rc *RunContext) context.Context {
	return context.WithValue(ctx, runContextKey, rc)
}
