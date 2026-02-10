package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"sort"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
	"golang.org/x/term"
)

// Formatter defines how to render output.
type Formatter interface {
	Format(data any) error
}

// Writers holds the stdout and stderr writers.
type Writers struct {
	Out io.Writer
	Err io.Writer
}

func DefaultWriters() *Writers {
	return &Writers{Out: os.Stdout, Err: os.Stderr}
}

// IsTTY returns true if the given file descriptor is a terminal.
func IsTTY(f *os.File) bool {
	return term.IsTerminal(int(f.Fd()))
}

// JSONFormatter outputs data as JSON.
type JSONFormatter struct {
	w io.Writer
}

func NewJSONFormatter(w io.Writer) *JSONFormatter {
	return &JSONFormatter{w: w}
}

func (f *JSONFormatter) Format(data any) error {
	enc := json.NewEncoder(f.w)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}

// TableFormatter outputs data as an ASCII table.
type TableFormatter struct {
	w io.Writer
}

func NewTableFormatter(w io.Writer) *TableFormatter {
	return &TableFormatter{w: w}
}

func (f *TableFormatter) Format(data any) error {
	switch v := data.(type) {
	case map[string]string:
		return f.formatMapStringString(v)
	case map[string]any:
		return f.formatMapStringAny(v)
	default:
		return f.formatReflect(data)
	}
}

func (f *TableFormatter) formatMapStringString(m map[string]string) error {
	table := newKeyValueTable(f.w)
	keys := sortedKeys(m)
	for _, k := range keys {
		_ = table.Append([]string{k, m[k]})
	}
	return table.Render()
}

func (f *TableFormatter) formatMapStringAny(m map[string]any) error {
	table := newKeyValueTable(f.w)
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		_ = table.Append([]string{k, formatValue(m[k])})
	}
	return table.Render()
}

func (f *TableFormatter) formatReflect(data any) error {
	rv := reflect.ValueOf(data)
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return nil
		}
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.Struct:
		return f.formatStruct(rv)
	case reflect.Slice:
		return f.formatSlice(rv)
	default:
		return NewJSONFormatter(f.w).Format(data)
	}
}

func (f *TableFormatter) formatStruct(rv reflect.Value) error {
	// Detect PaginatedResult[T]: has Items (slice), NextCursor, HasMore
	itemsField := rv.FieldByName("Items")
	nextCursorField := rv.FieldByName("NextCursor")
	hasMoreField := rv.FieldByName("HasMore")
	if itemsField.IsValid() && itemsField.Kind() == reflect.Slice &&
		nextCursorField.IsValid() && hasMoreField.IsValid() {
		if err := f.formatSlice(itemsField); err != nil {
			return err
		}
		if hasMoreField.Bool() {
			_, _ = fmt.Fprintf(f.w, "\nMore results available. Next cursor: %s\n", nextCursorField.String())
		}
		return nil
	}

	// Detect SearchResult: has Matches (slice) and Total (int)
	matchesField := rv.FieldByName("Matches")
	totalField := rv.FieldByName("Total")
	if matchesField.IsValid() && matchesField.Kind() == reflect.Slice &&
		totalField.IsValid() {
		if err := f.formatSlice(matchesField); err != nil {
			return err
		}
		_, _ = fmt.Fprintf(f.w, "\nTotal: %d\n", totalField.Int())
		return nil
	}

	// Generic single struct: render as key-value
	return f.formatStructAsKeyValue(rv)
}

func (f *TableFormatter) formatStructAsKeyValue(rv reflect.Value) error {
	rt := rv.Type()
	table := newKeyValueTable(f.w)
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		if !field.IsExported() {
			continue
		}
		name, skip := fieldName(field)
		if skip {
			continue
		}
		_ = table.Append([]string{name, formatValue(rv.Field(i).Interface())})
	}
	return table.Render()
}

func (f *TableFormatter) formatSlice(rv reflect.Value) error {
	if rv.Len() == 0 {
		_, _ = fmt.Fprintln(f.w, "No items found.")
		return nil
	}

	elemType := rv.Type().Elem()
	if elemType.Kind() == reflect.Ptr {
		elemType = elemType.Elem()
	}
	if elemType.Kind() != reflect.Struct {
		return NewJSONFormatter(f.w).Format(rv.Interface())
	}

	headers := structHeaders(elemType)
	table := newMultiColumnTable(f.w, headers)

	for i := 0; i < rv.Len(); i++ {
		elem := rv.Index(i)
		if elem.Kind() == reflect.Ptr {
			elem = elem.Elem()
		}
		_ = table.Append(structValues(elem, elemType))
	}
	return table.Render()
}

// PrintError writes an error message to stderr.
func PrintError(w *Writers, format string, args ...any) {
	_, _ = fmt.Fprintf(w.Err, "Error: "+format+"\n", args...)
}

// --- table helpers ---

func newKeyValueTable(w io.Writer) *tablewriter.Table {
	table := tablewriter.NewTable(w, tableOpts()...)
	table.Header([]string{"FIELD", "VALUE"})
	return table
}

func newMultiColumnTable(w io.Writer, headers []string) *tablewriter.Table {
	table := tablewriter.NewTable(w, tableOpts()...)
	table.Header(headers)
	return table
}

func tableOpts() []tablewriter.Option {
	return []tablewriter.Option{
		tablewriter.WithHeaderAlignment(tw.AlignLeft),
		tablewriter.WithRowAlignment(tw.AlignLeft),
		tablewriter.WithHeaderAutoWrap(tw.WrapNone),
		tablewriter.WithRowAutoWrap(tw.WrapNone),
	}
}

// --- reflection helpers ---

func fieldName(f reflect.StructField) (string, bool) {
	tag := f.Tag.Get("json")
	if tag == "-" {
		return "", true
	}
	if tag == "" {
		return f.Name, false
	}
	name, _, _ := strings.Cut(tag, ",")
	if name == "" {
		return f.Name, false
	}
	return name, false
}

func structHeaders(rt reflect.Type) []string {
	var headers []string
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		if !field.IsExported() {
			continue
		}
		name, skip := fieldName(field)
		if skip {
			continue
		}
		headers = append(headers, name)
	}
	return headers
}

func structValues(rv reflect.Value, rt reflect.Type) []string {
	var vals []string
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		if !field.IsExported() {
			continue
		}
		_, skip := fieldName(field)
		if skip {
			continue
		}
		vals = append(vals, formatValue(rv.Field(i).Interface()))
	}
	return vals
}

func formatValue(v any) string {
	if v == nil {
		return ""
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Slice {
		parts := make([]string, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			parts[i] = fmt.Sprintf("%v", rv.Index(i).Interface())
		}
		return strings.Join(parts, ", ")
	}
	return fmt.Sprintf("%v", v)
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
