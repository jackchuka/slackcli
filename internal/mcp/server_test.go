package mcp

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/jackchuka/slackcli/internal/slack/mocks"
	"github.com/jackchuka/slackcli/internal/version"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestNewServer_ReadWrite(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := mocks.NewMockService(ctrl)
	s := NewServer(mock, false)
	require.NotNil(t, s)
}

func TestNewServer_ReadOnly(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := mocks.NewMockService(ctrl)
	s := NewServer(mock, true)
	require.NotNil(t, s)
}

func TestToJSON(t *testing.T) {
	t.Run("struct", func(t *testing.T) {
		type data struct {
			Name string `json:"name"`
		}
		result := toJSON(data{Name: "test"})
		assert.Contains(t, result, `"name": "test"`)
	})

	t.Run("nil", func(t *testing.T) {
		result := toJSON(nil)
		assert.Equal(t, "null", result)
	})

	t.Run("map", func(t *testing.T) {
		result := toJSON(map[string]string{"key": "val"})
		assert.Contains(t, result, `"key": "val"`)
	})
}

func TestErrResult(t *testing.T) {
	result := errResult(errors.New("something broke"))
	require.NotNil(t, result)
	assert.True(t, result.IsError)
	require.NotEmpty(t, result.Content)
}

func TestNewServer_ReportsBuildVersion(t *testing.T) {
	orig := version.Version
	version.Version = "9.8.7"
	t.Cleanup(func() { version.Version = orig })

	ctrl := gomock.NewController(t)
	s := NewServer(mocks.NewMockService(ctrl), false)
	resp := s.HandleMessage(context.Background(), json.RawMessage(
		`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18","capabilities":{},"clientInfo":{"name":"test","version":"0"}}}`,
	))
	data, err := json.Marshal(resp)
	require.NoError(t, err)
	var got struct {
		Result struct {
			ServerInfo struct {
				Name    string `json:"name"`
				Version string `json:"version"`
			} `json:"serverInfo"`
		} `json:"result"`
	}
	require.NoError(t, json.Unmarshal(data, &got))
	assert.Equal(t, "slackcli", got.Result.ServerInfo.Name)
	assert.Equal(t, "9.8.7", got.Result.ServerInfo.Version)
}

func TestNewServer_ToolAnnotations(t *testing.T) {
	type hints struct{ readOnly, destructive bool }
	readOnly := hints{readOnly: true, destructive: false}
	additive := hints{readOnly: false, destructive: false}
	destructive := hints{readOnly: false, destructive: true}
	want := map[string]hints{
		"list_channels":       readOnly,
		"get_channel_info":    readOnly,
		"list_messages":       readOnly,
		"list_users":          readOnly,
		"get_user_info":       readOnly,
		"get_user_presence":   readOnly,
		"list_reactions":      readOnly,
		"list_files":          readOnly,
		"get_file_info":       readOnly,
		"search_messages":     readOnly,
		"auth_test":           readOnly,
		"create_channel":      additive,
		"invite_to_channel":   additive,
		"send_message":        additive,
		"add_reaction":        additive,
		"archive_channel":     destructive,
		"kick_from_channel":   destructive,
		"set_channel_topic":   destructive,
		"set_channel_purpose": destructive,
		"edit_message":        destructive,
		"delete_message":      destructive,
		"remove_reaction":     destructive,
		"delete_file":         destructive,
	}

	ctrl := gomock.NewController(t)
	tools := NewServer(mocks.NewMockService(ctrl), false).ListTools()
	require.Len(t, tools, len(want))
	for name, h := range want {
		tool, ok := tools[name]
		require.True(t, ok, name)
		a := tool.Tool.Annotations
		require.NotNil(t, a.ReadOnlyHint, name)
		require.NotNil(t, a.DestructiveHint, name)
		assert.Equal(t, h.readOnly, *a.ReadOnlyHint, name)
		assert.Equal(t, h.destructive, *a.DestructiveHint, name)
	}

	for name, tool := range NewServer(mocks.NewMockService(ctrl), true).ListTools() {
		assert.True(t, *tool.Tool.Annotations.ReadOnlyHint, name)
	}
}

func TestInviteToChannel_UserIDsItemSchema(t *testing.T) {
	ctrl := gomock.NewController(t)
	tool := NewServer(mocks.NewMockService(ctrl), false).GetTool("invite_to_channel")
	require.NotNil(t, tool)
	prop, ok := tool.Tool.InputSchema.Properties["user_ids"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, map[string]any{"type": "string"}, prop["items"])
}
