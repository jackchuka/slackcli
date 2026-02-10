package cmdutil

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunContext_SetGet(t *testing.T) {
	rc := &RunContext{ReadOnly: true}
	ctx := SetRunContext(context.Background(), rc)
	got := GetRunContext(ctx)

	require.NotNil(t, got)
	assert.Same(t, rc, got)
	assert.True(t, got.ReadOnly)
}

func TestGetRunContext_Panics(t *testing.T) {
	assert.Panics(t, func() {
		GetRunContext(context.Background())
	})
}
