package mcp

import (
	"errors"
	"testing"

	"github.com/jackchuka/slackcli/internal/slack/mocks"
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
