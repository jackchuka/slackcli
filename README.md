# slackcli

[![Test](https://github.com/jackchuka/slackcli/workflows/Test/badge.svg)](https://github.com/jackchuka/slackcli/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/jackchuka/slackcli)](https://goreportcard.com/report/github.com/jackchuka/slackcli)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**For AI, by AI.** The Slack CLI and MCP server built from the ground up for LLM-powered agents and programmatic access.

## Why slackcli?

Most Slack tools are built for humans. slackcli is built for agents. Every design decision — JSON-first output, structured error codes, automatic pagination, rate limit retries — is made so LLMs can interact with Slack reliably without hand-holding. Plug it in as an MCP server and your AI agent gets Slack superpowers. Use `--read-only` mode when you want guardrails.

## Quick Start

1. [Create a Slack app](https://api.slack.com/apps) and install it to your workspace with the scopes you need (e.g. `channels:read`, `chat:write`, `search:read`).
2. Grab the Bot or User OAuth Token (`xoxb-...`) from **OAuth & Permissions**.
3. Run:

```bash
brew install jackchuka/tap/slackcli
export SLACK_TOKEN=xoxb-your-token
slackcli channels list
```

## Features

- **CLI commands** for channels, messages, users, files, reactions, and search
- **MCP server** (stdio transport) for AI agent integration
- **JSON-first output** optimized for LLM consumption, with TTY-aware table fallback
- **Rate limit handling** with automatic retry
- **Error classification** with structured error codes
- **Pagination** support across all list operations
- **Read-only mode** to prevent accidental writes by AI agents

## Installation

Homebrew:

```bash
brew install jackchuka/tap/slackcli
```

Or via Go:

```bash
go install github.com/jackchuka/slackcli/cmd/slackcli@latest
```

Or build from source:

```bash
make build
```

## Authentication

Set your Slack token via environment variable or the CLI:

```bash
# Environment variable
export SLACK_TOKEN=xoxb-your-token

# Or use the auth command
slackcli auth login --token xoxb-your-token

# With a named workspace
slackcli auth login --token xoxb-your-token --name my-workspace
```

Token resolution order: `--token` flag > `SLACK_TOKEN` env > stored config.

## Usage

### CLI

```bash
# Channels
slackcli channels list
slackcli channels info C1234567890
slackcli channels create new-channel

# Messages
slackcli messages list --channel C1234567890
slackcli messages send --channel C1234567890 --text "Hello"
slackcli messages search --query "important"

# Users
slackcli users list
slackcli users info U1234567890

# Files
slackcli files list
slackcli files upload --channel C1234567890 --file ./report.pdf

# Reactions
slackcli reactions add --channel C1234567890 --timestamp 1234567890.123456 --name thumbsup
slackcli reactions list --user U1234567890
```

### MCP Server

Start the MCP server for AI agent integration:

```bash
slackcli mcp serve
```

Configure in your MCP client (e.g., Claude Desktop):

```json
{
  "mcpServers": {
    "slack": {
      "command": "slackcli",
      "args": ["mcp", "serve"],
      "env": {
        "SLACK_TOKEN": "xoxb-your-token"
      }
    }
  }
}
```

Available MCP tools:

| Category  | Tools                                                                                                                                                          |
| --------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Channels  | `list_channels`, `get_channel_info`, `create_channel`, `archive_channel`, `invite_to_channel`, `kick_from_channel`, `set_channel_topic`, `set_channel_purpose` |
| Messages  | `list_messages`, `send_message`, `edit_message`, `delete_message`, `search_messages`                                                                           |
| Users     | `list_users`, `get_user_info`, `get_user_presence`                                                                                                             |
| Reactions | `add_reaction`, `remove_reaction`, `list_reactions`                                                                                                            |
| Files     | `list_files`, `get_file_info`, `delete_file`                                                                                                                   |
| Auth      | `auth_test`                                                                                                                                                    |

### Read-Only Mode

Use `--read-only` to restrict to read-only operations. This prevents AI agents from accidentally sending messages, deleting files, or modifying channels.

**CLI** -- write commands are rejected before execution:

```bash
slackcli --read-only channels list           # works
slackcli --read-only messages send ...       # blocked
# Error: command "slackcli messages send" is a write operation and cannot be used in read-only mode
```

**MCP server** -- write tools are hidden entirely (AI agents never see them):

```json
{
  "mcpServers": {
    "slack-readonly": {
      "command": "slackcli",
      "args": ["--read-only", "mcp", "serve"],
      "env": {
        "SLACK_TOKEN": "xoxb-your-token"
      }
    }
  }
}
```

Read-only tools (always available): `auth_test`, `list_channels`, `get_channel_info`, `list_messages`, `list_users`, `get_user_info`, `get_user_presence`, `list_reactions`, `list_files`, `get_file_info`, `search_messages`.

Write tools (hidden in read-only mode): `create_channel`, `archive_channel`, `invite_to_channel`, `kick_from_channel`, `set_channel_topic`, `set_channel_purpose`, `send_message`, `edit_message`, `delete_message`, `add_reaction`, `remove_reaction`, `delete_file`.

## Output Formats

```bash
# Force JSON output
slackcli channels list -o json

# Force table output
slackcli channels list -o table
```

Default: table for TTY, JSON for piped output.

## License

MIT
