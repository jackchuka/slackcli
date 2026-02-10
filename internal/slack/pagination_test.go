package slack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPaginationParams_EffectiveLimit(t *testing.T) {
	tests := []struct {
		name  string
		limit int
		want  int
	}{
		{"zero returns default 100", 0, 100},
		{"negative returns default 100", -1, 100},
		{"positive returns value", 50, 50},
		{"one returns 1", 1, 1},
		{"large value returned as-is", 1000, 1000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := PaginationParams{Limit: tt.limit}
			assert.Equal(t, tt.want, p.EffectiveLimit())
		})
	}
}
