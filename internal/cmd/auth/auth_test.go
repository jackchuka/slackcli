package auth

import (
	"bytes"
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jackchuka/slackcli/internal/cmdutil"
	"github.com/jackchuka/slackcli/internal/config"
	"github.com/jackchuka/slackcli/internal/output"
)

func TestListCmd_NoWorkspacesPrintsEmptyArray(t *testing.T) {
	cfg, err := config.Load(filepath.Join(t.TempDir(), "config.json"))
	require.NoError(t, err)
	var buf bytes.Buffer
	rc := &cmdutil.RunContext{Config: cfg, Formatter: output.NewJSONFormatter(&buf)}

	cmd := newListCmd()
	cmd.SetContext(cmdutil.SetRunContext(context.Background(), rc))
	require.NoError(t, cmd.RunE(cmd, nil))
	assert.JSONEq(t, "[]", buf.String())
}
