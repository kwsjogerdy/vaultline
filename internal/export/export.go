// Package export provides functionality to export secrets in multiple formats.
package export

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
)

// Format represents an output format for secrets.
type Format string

const (
	FormatEnv  Format = "env"
	FormatJSON Format = "json"
	FormatCSV  Format = "csv"
)

// Exporter writes secrets in a given format.
type Exporter struct {
	format Format
	w      io.Writer
}

// New creates a new Exporter writing to w in the given format.
func New(format Format, w io.Writer) *Exporter {
	return &Exporter{format: format, w: w}
}

// Write outputs the secrets map in the configured format.
func (e *Exporter) Write(secrets map[string]string) error {
	switch e.format {
	case FormatEnv:
		return e.writeEnv(secrets)
	case FormatJSON:
		return e.writeJSON(secrets)
	case FormatCSV:
		return e.writeCSV(secrets)
	default:
		return fmt.Errorf("unsupported format: %s", e.format)
	}
}

func (e *Exporter) writeEnv(secrets map[string]string) error {
	for _, k := range sortedKeys(secrets) {
		v := secrets[k]
		if strings.ContainsAny(v, " \t\n") {
			v = fmt.Sprintf("%q", v)
		}
		if _, err := fmt.Fprintf(e.w, "%s=%s\n", k, v); err != nil {
			return err
		}
	}
	return nil
}

func (e *Exporter) writeJSON(secrets map[string]string) error {
	enc := json.NewEncoder(e.w)
	enc.SetIndent("", "  ")
	return enc.Encode(secrets)
}

func (e *Exporter) writeCSV(secrets map[string]string) error {
	if _, err := fmt.Fprintln(e.w, "key,value"); err != nil {
		return err
	}
	for _, k := range sortedKeys(secrets) {
		v := strings.ReplaceAll(secrets[k], "\"", "\"\"")
		if _, err := fmt.Fprintf(e.w, "%s,\"%s\"\n", k, v); err != nil {
			return err
		}
	}
	return nil
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
