# Regenerate with: vhs assets/demo.tape (from the repository root)
# Runs fully offline: no Slack API call is made, and a placeholder token keeps real credentials out.

Output assets/demo.gif

Set Shell "bash"
Set FontSize 22
Set Width 1440
Set Height 760
Set Padding 24
Set Theme "Catppuccin Mocha"
Set TypingSpeed 45ms

Hide
Type `export D=/tmp/slackcli-demo && rm -rf "$D" && mkdir -p "$D/home" && go build -o "$D/slackcli" ./cmd/slackcli && cp assets/demo/list-tools.jsonl "$D" && cd "$D" && export HOME="$D/home" XDG_CONFIG_HOME="$D/home/.config" SLACK_TOKEN=xoxb-placeholder PATH="$D:$PATH" PS1='\[\e[1;35m\]$ \[\e[0m\]' && clear`
Enter
Sleep 4s
Show

Type "slackcli --help"
Sleep 400ms
Enter
Sleep 4s

Type "clear"
Enter
Type "slackcli mcp serve < list-tools.jsonl | jq '.result.tools // empty | length'"
Sleep 400ms
Enter
Sleep 3s

Type "slackcli --read-only mcp serve < list-tools.jsonl | jq -r '.result.tools[]?.name'"
Sleep 400ms
Enter
Sleep 4s

Type "clear"
Enter
Type `slackcli --read-only messages send --channel C0123456789 --text "hello"`
Sleep 400ms
Enter
Sleep 4s
