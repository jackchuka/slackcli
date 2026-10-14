<div align="center">

# slackcli

**A Slack CLI and MCP server built for AI agents.**

JSON-first output, structured errors, and automatic pagination,<br>
with a read-only mode that keeps agents from posting or deleting.

[![Test](https://github.com/jackchuka/slackcli/actions/workflows/test.yml/badge.svg)](https://github.com/jackchuka/slackcli/actions/workflows/test.yml)
[![Release](https://img.shields.io/github/v/release/jackchuka/slackcli?sort=semver)](https://github.com/jackchuka/slackcli/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

<img src="assets/demo.gif" alt="slackcli listing its MCP tools, hiding write tools in read-only mode, and blocking a write command" width="800">

[Installation](#installation) · [Quick Start](#quick-start) · [Commands](#commands) · [MCP Server](#mcp-server) · [Read-Only Mode](#read-only-mode)

</div>

---

Most Slack tools are built for people at a keyboard. An agent needs something else: output it can parse, errors it can branch on, pagination it doesn't have to manage, and a guarantee that it can't post to `#general` by accident. **slackcli** gives you that as a single binary. Run it from scripts, or plug it into any MCP client as a stdio server.

## Features

- **CLI and MCP in one binary**: channels, messages, users, files, reactions and search, from the shell or as MCP tools
- **JSON-first output**: JSON when piped, a table in a terminal, or forced with `-o json|table`
- **Read-only mode**: `--read-only` rejects write commands in the CLI and hides write tools from MCP clients
- **Automatic pagination**: cursor-based paging on list commands, plus `--all` to fetch every page
- **Rate-limit retries**: rate-limited calls are retried up to 3 times, honoring Slack's `Retry-After`; a call still limited after that fails with `rate_limited`
- **Structured errors**: Slack errors are classified (`auth_error`, `not_found`, `rate_limited`, `permission_denied`, ...) and mapped to exit codes
- **Multiple workspaces**: store several tokens, switch between them with `auth switch`, or pick one per command with `-w`

## Installation

### Homebrew

```bash
brew install jackchuka/tap/slackcli
```

### Go

```bash
go install github.com/jackchuka/slackcli/cmd/slackcli@latest
```

### Binary

Download a prebuilt binary for Linux, macOS, or Windows from the [releases page](https://github.com/jackchuka/slackcli/releases).

### From Source

```bash
git clone https://github.com/jackchuka/slackcli.git
cd slackcli
make build   # writes bin/slackcli
```

## Quick Start

1. **Create a Slack app** at [api.slack.com/apps](https://api.slack.com/apps), add the scopes you need (for example `channels:read`, `chat:write`, `search:read`), and install it to your workspace.

2. **Copy the token** (`xoxb-...` or `xoxp-...`) from **OAuth & Permissions**, and either export it or save it:

   ```bash
   export SLACK_TOKEN=xoxb-your-token
   # or store it in the config file
   slackcli auth login --token xoxb-your-token
   ```

3. **Run a command**:

   ```bash
   slackcli channels list
   slackcli messages search --query "deploy" --limit 5
   ```

## Authentication

The token is resolved in this order:

1. `--token` flag
2. `-w, --workspace` flag: a workspace saved with `auth login`
3. `SLACK_TOKEN` environment variable
4. The active workspace in the config file

`auth login` verifies the token with Slack, then saves it under a workspace name (the team name by default, or `--name`). The config file lives at `slackcli/config.json` in your OS config directory: `~/.config` on Linux, `~/Library/Application Support` on macOS.

```bash
slackcli auth login --token xoxb-your-token --name team-a
slackcli auth list            # configured workspaces
slackcli auth switch team-a   # change the active workspace
slackcli auth status          # check the current token against Slack
slackcli auth logout [name]   # remove a workspace (default: the active one)
```

`-w` selects a saved workspace for a single command without changing the active one. `auth status` and `auth logout` act on it too:

```bash
slackcli -w team-b channels list
slackcli -w team-b auth status
```

An unknown workspace name is an error.

## Commands

Arguments in `<angle brackets>` are positional. Flags marked * are required.

### Channels

| Command                                     | Description                                        |
| ------------------------------------------- | -------------------------------------------------- |
| `channels list`                             | List channels (`--limit` 100, `--cursor`, `--all`) |
| `channels info <channel-id>`                | Get channel details                                |
| `channels create <name>`                    | Create a channel (`--private` for a private one)   |
| `channels archive <channel-id>`             | Archive a channel                                  |
| `channels invite <channel-id> <user-id>...` | Invite one or more users                           |
| `channels kick <channel-id> <user-id>`      | Remove a user                                      |
| `channels topic <channel-id> <topic>`       | Set the topic                                      |
| `channels purpose <channel-id> <purpose>`   | Set the purpose                                    |

### Messages

| Command           | Description                                                                                                                         |
| ----------------- | ----------------------------------------------------------------------------------------------------------------------------------- |
| `messages list`   | Channel history: `--channel`\*, `--limit` 100, `--cursor`, `--all`                                                                  |
| `messages send`   | Post a message: `--channel`\*, `--text`\*                                                                                           |
| `messages reply`  | Reply in a thread: `--channel`\*, `--thread-ts`\*, `--text`\*                                                                       |
| `messages edit`   | Edit a message: `--channel`\*, `--timestamp`\*, `--text`\*                                                                          |
| `messages delete` | Delete a message: `--channel`\*, `--timestamp`\*                                                                                    |
| `messages search` | Search: `--query`\*, `--sort` `timestamp`\|`score` (default `timestamp`), `--sort-dir` `asc`\|`desc` (default `desc`), `--limit` 20 |

### Users

| Command                    | Description                         |
| -------------------------- | ----------------------------------- |
| `users list`               | List users (`--limit` 100, `--all`) |
| `users info <user-id>`     | Get a user profile                  |
| `users presence <user-id>` | Get presence (`active` or `away`)   |

### Files

| Command                    | Description                                                                         |
| -------------------------- | ----------------------------------------------------------------------------------- |
| `files list`               | List files (`--channel`, `--user`, `--limit` 100, `--cursor`, `--all`)              |
| `files info <file-id>`     | Get file metadata                                                                   |
| `files upload`             | Upload a file: `--channel`\*, `--file`\*, `--title`                                 |
| `files download <file-id>` | Download a file to `-d, --dest` (default: the file's name in the current directory) |
| `files delete <file-id>`   | Delete a file                                                                       |

### Reactions

| Command            | Description                                                                                |
| ------------------ | ------------------------------------------------------------------------------------------ |
| `reactions add`    | Add a reaction: `--channel`\*, `--timestamp`\*, `--name`\* (emoji name, no colons)         |
| `reactions remove` | Remove a reaction: same flags as `add`                                                     |
| `reactions list`   | List items a user reacted to (`--user`, default: the token's user; `--limit` 100, `--all`) |

### Other

| Command      | Description                          |
| ------------ | ------------------------------------ |
| `mcp serve`  | Start the MCP server over stdio      |
| `version`    | Print version, commit and build date |
| `completion` | Generate shell completion scripts    |

### Global flags

| Flag                     | Description                                                  |
| ------------------------ | ------------------------------------------------------------ |
| `--token <token>`        | Slack API token; overrides `SLACK_TOKEN` and the config      |
| `-w, --workspace <name>` | Use this saved workspace instead of the active one           |
| `-o, --output <fmt>`     | `json` or `table`; default is table in a TTY, JSON otherwise |
| `--read-only`            | Reject write commands, and hide write tools in `mcp serve`   |

### Examples

```bash
slackcli channels list --all -o json | jq -r '.items[].name'
slackcli messages send --channel C0123456789 --text "Build finished"
slackcli messages reply --channel C0123456789 --thread-ts 1700000000.000100 --text "Done"
slackcli files upload --channel C0123456789 --file ./report.pdf --title "Weekly report"
slackcli reactions add --channel C0123456789 --timestamp 1700000000.000100 --name thumbsup
```

## MCP Server

`slackcli mcp serve` speaks MCP over stdio. Add it to your MCP client's config:

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

If `SLACK_TOKEN` is not set, the server uses the active workspace from `slackcli auth login`. To serve a different saved workspace, use `"args": ["-w", "team-b", "mcp", "serve"]`. The server reports slackcli's version to the client.

| Category  | Read tools                                         | Write tools                                                                                                               |
| --------- | -------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------- |
| Channels  | `list_channels`, `get_channel_info`                | `create_channel`, `archive_channel`, `invite_to_channel`, `kick_from_channel`, `set_channel_topic`, `set_channel_purpose` |
| Messages  | `list_messages`, `search_messages`                 | `send_message`, `edit_message`, `delete_message`                                                                          |
| Users     | `list_users`, `get_user_info`, `get_user_presence` |                                                                                                                           |
| Reactions | `list_reactions`                                   | `add_reaction`, `remove_reaction`                                                                                         |
| Files     | `list_files`, `get_file_info`                      | `delete_file`                                                                                                             |
| Auth      | `auth_test`                                        |                                                                                                                           |

## Read-Only Mode

Add `--read-only` when an agent should be able to look but not touch.

In the **CLI**, write commands fail before any API call:

```console
$ slackcli --read-only messages send --channel C0123456789 --text "hello"
Error: command "slackcli messages send" is a write operation and cannot be used in read-only mode
```

In the **MCP server**, write tools aren't registered at all, so the agent never sees them. Only the 11 read tools from the table above are exposed:

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

## Output and Exit Codes

List commands return `{"items": [...], "next_cursor": "...", "has_more": true}`. Pass `next_cursor` back with `--cursor`, or use `--all` to fetch every page. In table mode, the next cursor is printed below the table.

| Exit code | Meaning                                     |
| --------- | ------------------------------------------- |
| `0`       | Success                                     |
| `1`       | Any other error                             |
| `2`       | Authentication error, or no token found     |
| `3`       | Not found (channel, user, file, or message) |
| `4`       | Still rate limited after 3 retries          |

## License

[MIT](LICENSE)
