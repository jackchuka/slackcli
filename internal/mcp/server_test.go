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
