package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	t.Run("non-existent file returns empty config", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "nonexistent.json")
		cfg, err := Load(path)

		require.NoError(t, err)
		assert.NotNil(t, cfg.Workspaces)
		assert.Empty(t, cfg.Workspaces)
		assert.Empty(t, cfg.ActiveWorkspace)
	})

	t.Run("valid JSON loads correctly", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "config.json")
		data := `{
			"active_workspace": "team1",
			"workspaces": {
				"team1": {"name": "Team One", "token": "xoxb-123", "team_id": "T1"}
			}
		}`
		require.NoError(t, os.WriteFile(path, []byte(data), 0o600))

		cfg, err := Load(path)

		require.NoError(t, err)
		assert.Equal(t, "team1", cfg.ActiveWorkspace)
		require.Contains(t, cfg.Workspaces, "team1")
		assert.Equal(t, "Team One", cfg.Workspaces["team1"].Name)
		assert.Equal(t, "xoxb-123", cfg.Workspaces["team1"].Token)
		assert.Equal(t, "T1", cfg.Workspaces["team1"].TeamID)
	})

	t.Run("invalid JSON returns error", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "config.json")
		require.NoError(t, os.WriteFile(path, []byte("{invalid"), 0o600))

		_, err := Load(path)

		assert.Error(t, err)
	})

	t.Run("null workspaces initializes to empty map", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "config.json")
		data := `{"active_workspace": "", "workspaces": null}`
		require.NoError(t, os.WriteFile(path, []byte(data), 0o600))

		cfg, err := Load(path)

		require.NoError(t, err)
		assert.NotNil(t, cfg.Workspaces)
		assert.Empty(t, cfg.Workspaces)
	})
}

func TestConfig_Save(t *testing.T) {
	t.Run("roundtrip Load -> modify -> Save -> Load", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "sub", "config.json")

		cfg, err := Load(path)
		require.NoError(t, err)

		cfg.SetWorkspace("myteam", Workspace{
			Name:   "My Team",
			Token:  "xoxb-test",
			TeamID: "T123",
		})

		require.NoError(t, cfg.Save())

		// Verify file was created with correct permissions
		info, err := os.Stat(path)
		require.NoError(t, err)
		assert.Equal(t, os.FileMode(0o600), info.Mode().Perm())

		// Reload and verify
		cfg2, err := Load(path)
		require.NoError(t, err)
		assert.Equal(t, "myteam", cfg2.ActiveWorkspace)
		require.Contains(t, cfg2.Workspaces, "myteam")
		assert.Equal(t, "xoxb-test", cfg2.Workspaces["myteam"].Token)
	})
}

func TestConfig_ActiveToken(t *testing.T) {
	tests := []struct {
		name            string
		activeWorkspace string
		workspaces      map[string]Workspace
		want            string
	}{
		{
			name:            "returns token for active workspace",
			activeWorkspace: "team1",
			workspaces:      map[string]Workspace{"team1": {Token: "xoxb-123"}},
			want:            "xoxb-123",
		},
		{
			name:            "returns empty when no active workspace",
			activeWorkspace: "",
			workspaces:      map[string]Workspace{"team1": {Token: "xoxb-123"}},
			want:            "",
		},
		{
			name:            "returns empty when active workspace not in map",
			activeWorkspace: "missing",
			workspaces:      map[string]Workspace{},
			want:            "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				ActiveWorkspace: tt.activeWorkspace,
				Workspaces:      tt.workspaces,
			}
			assert.Equal(t, tt.want, cfg.ActiveToken())
		})
	}
}

func TestConfig_SetWorkspace(t *testing.T) {
	cfg := &Config{Workspaces: make(map[string]Workspace)}

	cfg.SetWorkspace("team1", Workspace{Name: "Team One", Token: "xoxb-1"})

	assert.Equal(t, "team1", cfg.ActiveWorkspace)
	require.Contains(t, cfg.Workspaces, "team1")
	assert.Equal(t, "xoxb-1", cfg.Workspaces["team1"].Token)

	// Adding a second workspace switches active
	cfg.SetWorkspace("team2", Workspace{Name: "Team Two", Token: "xoxb-2"})

	assert.Equal(t, "team2", cfg.ActiveWorkspace)
	assert.Len(t, cfg.Workspaces, 2)
}

func TestConfig_RemoveWorkspace(t *testing.T) {
	t.Run("removes active workspace and switches to another", func(t *testing.T) {
		cfg := &Config{
			ActiveWorkspace: "team1",
			Workspaces: map[string]Workspace{
				"team1": {Token: "xoxb-1"},
				"team2": {Token: "xoxb-2"},
			},
		}

		cfg.RemoveWorkspace("team1")

		assert.NotContains(t, cfg.Workspaces, "team1")
		assert.Equal(t, "team2", cfg.ActiveWorkspace)
	})

	t.Run("removes only workspace", func(t *testing.T) {
		cfg := &Config{
			ActiveWorkspace: "team1",
			Workspaces: map[string]Workspace{
				"team1": {Token: "xoxb-1"},
			},
		}

		cfg.RemoveWorkspace("team1")

		assert.Empty(t, cfg.Workspaces)
		assert.Empty(t, cfg.ActiveWorkspace)
	})

	t.Run("removes non-active workspace", func(t *testing.T) {
		cfg := &Config{
			ActiveWorkspace: "team1",
			Workspaces: map[string]Workspace{
				"team1": {Token: "xoxb-1"},
				"team2": {Token: "xoxb-2"},
			},
		}

		cfg.RemoveWorkspace("team2")

		assert.NotContains(t, cfg.Workspaces, "team2")
		assert.Equal(t, "team1", cfg.ActiveWorkspace)
	})

	t.Run("removing non-existent workspace is no-op", func(t *testing.T) {
		cfg := &Config{
			ActiveWorkspace: "team1",
			Workspaces: map[string]Workspace{
				"team1": {Token: "xoxb-1"},
			},
		}

		cfg.RemoveWorkspace("nonexistent")

		assert.Len(t, cfg.Workspaces, 1)
		assert.Equal(t, "team1", cfg.ActiveWorkspace)
	})
}
