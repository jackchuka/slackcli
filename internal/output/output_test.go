package output

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONFormatter_Format(t *testing.T) {
	t.Run("formats map as indented JSON", func(t *testing.T) {
		var buf bytes.Buffer
		f := NewJSONFormatter(&buf)
		data := map[string]string{"key": "value"}

		require.NoError(t, f.Format(data))

		var result map[string]string
		require.NoError(t, json.Unmarshal(buf.Bytes(), &result))
		assert.Equal(t, "value", result["key"])
	})

	t.Run("formats struct", func(t *testing.T) {
		var buf bytes.Buffer
		f := NewJSONFormatter(&buf)
		type item struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		}

		require.NoError(t, f.Format(item{ID: "1", Name: "test"}))

		assert.Contains(t, buf.String(), `"id": "1"`)
		assert.Contains(t, buf.String(), `"name": "test"`)
	})

	t.Run("output is indented", func(t *testing.T) {
		var buf bytes.Buffer
		f := NewJSONFormatter(&buf)

		require.NoError(t, f.Format(map[string]int{"a": 1}))

		assert.Contains(t, buf.String(), "  ")
	})
}

func TestTableFormatter_Format(t *testing.T) {
	t.Run("map[string]string renders key-value table", func(t *testing.T) {
		var buf bytes.Buffer
		f := NewTableFormatter(&buf)
		data := map[string]string{"name": "general", "id": "C123"}

		require.NoError(t, f.Format(data))

		out := buf.String()
		assert.Contains(t, out, "id")
		assert.Contains(t, out, "C123")
		assert.Contains(t, out, "name")
		assert.Contains(t, out, "general")
	})

	t.Run("map[string]any renders key-value table", func(t *testing.T) {
		var buf bytes.Buffer
		f := NewTableFormatter(&buf)
		data := map[string]any{"count": 42, "name": "test"}

		require.NoError(t, f.Format(data))

		out := buf.String()
		assert.Contains(t, out, "count")
		assert.Contains(t, out, "42")
	})

	t.Run("nil pointer renders nothing", func(t *testing.T) {
		var buf bytes.Buffer
		f := NewTableFormatter(&buf)
		var p *struct{ Name string }

		require.NoError(t, f.Format(p))

		assert.Empty(t, buf.String())
	})

	t.Run("PaginatedResult shape with HasMore", func(t *testing.T) {
		type Item struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		}
		type Result struct {
			Items      []Item `json:"items"`
			NextCursor string `json:"next_cursor"`
			HasMore    bool   `json:"has_more"`
		}

		var buf bytes.Buffer
		f := NewTableFormatter(&buf)
		data := Result{
			Items:      []Item{{ID: "1", Name: "first"}, {ID: "2", Name: "second"}},
			NextCursor: "cursor_abc",
			HasMore:    true,
		}

		require.NoError(t, f.Format(data))

		out := buf.String()
		assert.Contains(t, out, "first")
		assert.Contains(t, out, "second")
		assert.Contains(t, out, "More results available")
		assert.Contains(t, out, "cursor_abc")
	})

	t.Run("PaginatedResult shape without HasMore", func(t *testing.T) {
		type Item struct {
			ID string `json:"id"`
		}
		type Result struct {
			Items      []Item `json:"items"`
			NextCursor string `json:"next_cursor"`
			HasMore    bool   `json:"has_more"`
		}

		var buf bytes.Buffer
		f := NewTableFormatter(&buf)
		data := Result{
			Items:   []Item{{ID: "1"}},
			HasMore: false,
		}

		require.NoError(t, f.Format(data))

		out := buf.String()
		assert.Contains(t, out, "1")
		assert.NotContains(t, out, "More results available")
	})

	t.Run("SearchResult shape", func(t *testing.T) {
		type Match struct {
			Text string `json:"text"`
		}
		type SearchResult struct {
			Matches []Match `json:"matches"`
			Total   int     `json:"total"`
		}

		var buf bytes.Buffer
		f := NewTableFormatter(&buf)
		data := SearchResult{
			Matches: []Match{{Text: "hello world"}},
			Total:   42,
		}

		require.NoError(t, f.Format(data))

		out := buf.String()
		assert.Contains(t, out, "hello world")
		assert.Contains(t, out, "Total: 42")
	})

	t.Run("single struct renders as key-value", func(t *testing.T) {
		type Info struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}

		var buf bytes.Buffer
		f := NewTableFormatter(&buf)

		require.NoError(t, f.Format(Info{Name: "Alice", Age: 30}))

		out := buf.String()
		assert.Contains(t, out, "name")
		assert.Contains(t, out, "Alice")
		assert.Contains(t, out, "age")
		assert.Contains(t, out, "30")
	})

	t.Run("empty slice renders no items message", func(t *testing.T) {
		type Item struct {
			ID string `json:"id"`
		}
		type Result struct {
			Items      []Item `json:"items"`
			NextCursor string `json:"next_cursor"`
			HasMore    bool   `json:"has_more"`
		}

		var buf bytes.Buffer
		f := NewTableFormatter(&buf)

		require.NoError(t, f.Format(Result{Items: []Item{}}))

		assert.Contains(t, buf.String(), "No items found.")
	})

	t.Run("slice of structs renders multi-column table", func(t *testing.T) {
		type Row struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		}

		var buf bytes.Buffer
		f := NewTableFormatter(&buf)

		require.NoError(t, f.Format([]Row{
			{ID: "1", Name: "alpha"},
			{ID: "2", Name: "beta"},
		}))

		out := buf.String()
		// tablewriter uppercases headers
		assert.Contains(t, out, "ID")
		assert.Contains(t, out, "NAME")
		assert.Contains(t, out, "alpha")
		assert.Contains(t, out, "beta")
	})

	t.Run("json tag '-' skips field", func(t *testing.T) {
		type S struct {
			Visible string `json:"visible"`
			Hidden  string `json:"-"`
		}

		var buf bytes.Buffer
		f := NewTableFormatter(&buf)

		require.NoError(t, f.Format(S{Visible: "yes", Hidden: "no"}))

		out := buf.String()
		assert.Contains(t, out, "visible")
		assert.Contains(t, out, "yes")
		assert.NotContains(t, out, "Hidden")
		assert.NotContains(t, out, "no")
	})
}

func TestPrintError(t *testing.T) {
	var errBuf bytes.Buffer
	w := &Writers{Out: &bytes.Buffer{}, Err: &errBuf}

	PrintError(w, "something went %s", "wrong")

	assert.Equal(t, "Error: something went wrong\n", errBuf.String())
}

func TestFieldName(t *testing.T) {
	tests := []struct {
		name     string
		tag      string
		wantName string
		wantSkip bool
	}{
		{"json dash skips", "-", "", true},
		{"custom name", "custom_name", "custom_name", false},
		{"name with omitempty", "custom_name,omitempty", "custom_name", false},
		{"empty tag uses field name", "", "", false},
		{"only omitempty uses field name", ",omitempty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := reflect.StructField{
				Name: "FieldName",
				Tag:  reflect.StructTag(`json:"` + tt.tag + `"`),
			}
			name, skip := fieldName(field)

			assert.Equal(t, tt.wantSkip, skip)
			if !skip {
				if tt.wantName == "" {
					assert.Equal(t, "FieldName", name)
				} else {
					assert.Equal(t, tt.wantName, name)
				}
			}
		})
	}

	t.Run("no json tag uses field name", func(t *testing.T) {
		field := reflect.StructField{Name: "MyField"}
		name, skip := fieldName(field)
		assert.False(t, skip)
		assert.Equal(t, "MyField", name)
	})
}

func TestFormatValue(t *testing.T) {
	tests := []struct {
		name  string
		input any
		want  string
	}{
		{"nil", nil, ""},
		{"string", "hello", "hello"},
		{"int", 42, "42"},
		{"bool true", true, "true"},
		{"bool false", false, "false"},
		{"string slice", []string{"a", "b", "c"}, "a, b, c"},
		{"int slice", []int{1, 2, 3}, "1, 2, 3"},
		{"empty slice", []string{}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, formatValue(tt.input))
		})
	}
}

func TestSortedKeys(t *testing.T) {
	t.Run("returns sorted keys", func(t *testing.T) {
		m := map[string]string{"zebra": "z", "apple": "a", "mango": "m"}
		got := sortedKeys(m)
		assert.Equal(t, []string{"apple", "mango", "zebra"}, got)
	})

	t.Run("empty map returns empty slice", func(t *testing.T) {
		got := sortedKeys(map[string]string{})
		assert.Empty(t, got)
	})
}
