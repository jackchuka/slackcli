package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jackchuka/slackcli/internal/cmdutil"
	"github.com/jackchuka/slackcli/internal/config"
)

func writeTestConfig(t *testing.T) {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, ".config"))
	t.Setenv("SLACK_TOKEN", "")
	path := config.DefaultPath()
	require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o700))
	data := `{
		"active_workspace": "team-a",
		"workspaces": {
			"team-a": {"name": "team-a", "token": "xoxb-a"},
			"team-b": {"name": "team-b", "token": "xoxb-b"}
		}
	}`
	require.NoError(t, os.WriteFile(path, []byte(data), 0o600))
}

func runWithProbe(t *testing.T, args ...string) (*cmdutil.RunContext, error) {
	t.Helper()
	var rc *cmdutil.RunContext
	root := NewRootCmd()
	root.AddCommand(&cobra.Command{
		Use: "probe",
		RunE: func(c *cobra.Command, _ []string) error {
			rc = cmdutil.GetRunContext(c.Context())
			return nil
		},
	})
	root.SetArgs(args)
	err := root.Execute()
	return rc, err
}

func TestWorkspaceFlag(t *testing.T) {
	t.Run("selects the named workspace's token", func(t *testing.T) {
		writeTestConfig(t)
		rc, err := runWithProbe(t, "-w", "team-b", "probe")
		require.NoError(t, err)
		assert.Equal(t, "xoxb-b", rc.Resolver.Resolve())
		assert.Equal(t, "team-b", rc.Workspace)
	})

	t.Run("overrides SLACK_TOKEN", func(t *testing.T) {
		writeTestConfig(t)
		t.Setenv("SLACK_TOKEN", "xoxb-env")
		rc, err := runWithProbe(t, "--workspace", "team-b", "probe")
		require.NoError(t, err)
		assert.Equal(t, "xoxb-b", rc.Resolver.Resolve())
	})

	t.Run("--token still wins", func(t *testing.T) {
		writeTestConfig(t)
		rc, err := runWithProbe(t, "--token", "xoxb-flag", "-w", "team-b", "probe")
		require.NoError(t, err)
		assert.Equal(t, "xoxb-flag", rc.Resolver.Resolve())
	})

	t.Run("defaults to the active workspace", func(t *testing.T) {
		writeTestConfig(t)
		rc, err := runWithProbe(t, "probe")
		require.NoError(t, err)
		assert.Equal(t, "xoxb-a", rc.Resolver.Resolve())
		assert.Equal(t, "team-a", rc.Workspace)
	})

	t.Run("unknown workspace is an error", func(t *testing.T) {
		writeTestConfig(t)
		_, err := runWithProbe(t, "-w", "nope", "probe")
		require.Error(t, err)
		assert.Contains(t, err.Error(), `workspace "nope" not found`)
	})
}
