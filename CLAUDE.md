# SlackCLI

## Build & Test
- Build: `make build` or `go build ./cmd/slackcli`
- Test: `make test` or `go test ./... -race -count=1`
- Lint: `make lint` or `golangci-lint run ./...`
- Single test: `go test ./internal/slack/ -run TestName -v`

## Architecture
- CLI (Cobra) + MCP server share `internal/slack/` service layer
- `internal/cmd/` — Cobra command definitions
- `internal/slack/` — Slack API client with pagination, rate limiting, error classification
- `internal/config/` — Config management (XDG paths)
- `internal/auth/` — Token resolution chain: flag → env → config
- `internal/output/` — JSON/Table formatters, TTY detection
- `internal/mcp/` — MCP server and tool handlers

## Conventions
- JSON output by default (LLM-first)
- All list operations use `PaginationParams` / `PaginatedResult[T]`
- Errors are classified via `SlackError` with `ErrorCode`
- No global state — `RunContext` passed via Cobra context
- Exit codes: 0=ok, 1=error, 2=auth, 3=not found
