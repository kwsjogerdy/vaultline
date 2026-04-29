// Package envdoc generates human-readable documentation for secret maps.
package envdoc

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// Entry holds documentation metadata for a single secret key.
type Entry struct {
	Key         string
	Description string
	Example     string
	Required    bool
}

// Doc holds the full documentation set for a secret map.
type Doc struct {
	entries map[string]Entry
}

// New creates a new Doc populated from a map of annotations.
// annotations maps key -> "description|example|required" (pipe-separated).
func New(annotations map[string]string) (*Doc, error) {
	d := &Doc{entries: make(map[string]Entry)}
	for key, raw := range annotations {
		if key == "" {
			return nil, fmt.Errorf("envdoc: empty key")
		}
		parts := strings.SplitN(raw, "|", 3)
		e := Entry{Key: key}
		if len(parts) > 0 {
			e.Description = strings.TrimSpace(parts[0])
		}
		if len(parts) > 1 {
			e.Example = strings.TrimSpace(parts[1])
		}
		if len(parts) > 2 {
			e.Required = strings.TrimSpace(parts[2]) == "required"
		}
		d.entries[key] = e
	}
	return d, nil
}

// Get returns the Entry for a given key and whether it was found.
func (d *Doc) Get(key string) (Entry, bool) {
	e, ok := d.entries[key]
	return e, ok
}

// Keys returns all documented keys in sorted order.
func (d *Doc) Keys() []string {
	keys := make([]string, 0, len(d.entries))
	for k := range d.entries {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Write renders documentation to w in the given format ("text" or "markdown").
func (d *Doc) Write(w io.Writer, format string) error {
	keys := d.Keys()
	switch strings.ToLower(format) {
	case "markdown":
		fmt.Fprintln(w, "| Key | Description | Example | Required |")
		fmt.Fprintln(w, "|-----|-------------|---------|----------|")
		for _, k := range keys {
			e := d.entries[k]
			req := "no"
			if e.Required {
				req = "yes"
			}
			fmt.Fprintf(w, "| %s | %s | %s | %s |\n", e.Key, e.Description, e.Example, req)
		}
	case "text":
		for _, k := range keys {
			e := d.entries[k]
			req := ""
			if e.Required {
				req = " [required]"
			}
			fmt.Fprintf(w, "%s%s\n", e.Key, req)
			if e.Description != "" {
				fmt.Fprintf(w, "  Description: %s\n", e.Description)
			}
			if e.Example != "" {
				fmt.Fprintf(w, "  Example:     %s\n", e.Example)
			}
		}
	default:
		return fmt.Errorf("envdoc: unsupported format %q", format)
	}
	return nil
}
