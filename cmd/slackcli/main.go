package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/jackchuka/slackcli/internal/cmd"
	"github.com/jackchuka/slackcli/internal/slack"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(exitCode(err))
	}
}

func exitCode(err error) int {
	var se *slack.SlackError
	if errors.As(err, &se) {
		switch se.Code {
		case slack.ErrAuth:
			return 2
		case slack.ErrNotFound:
			return 3
		}
	}
	return 1
}
