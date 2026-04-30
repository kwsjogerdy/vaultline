// Package envformat normalises the textual representation of secret maps
// into a canonical .env line format, supporting optional sorting, comment
// headers, and blank-line separators between prefix groups.
package envformat

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// Options controls how the formatted output is rendered.
type Options struct {
	SortKeys   bool   // sort keys alphabetically
	GroupByPrefix bool // insert blank line between prefix groups
	Header     string // optional comment block prepended to output
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		SortKeys:      true,
		GroupByPrefix: false,
	}
}

// Formatter renders a secret map to a writer.
type Formatter struct {
	opts Options
}

// New creates a Formatter with the given options.
func New(opts Options) *Formatter {
	return &Formatter{opts: opts}
}

// Write renders secrets to w in .env format.
func (f *Formatter) Write(w io.Writer, secrets map[string]string) error {
	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	if f.opts.SortKeys {
		sort.Strings(keys)
	}

	if f.opts.Header != "" {
		for _, line := range strings.Split(f.opts.Header, "\n") {
			if _, err := fmt.Fprintf(w, "# %s\n", line); err != nil {
				return err
			}
		}
		if _, err := fmt.Fprintln(w); err != nil {
			return err
		}
	}

	lastPrefix := ""
	for _, k := range keys {
		if f.opts.GroupByPrefix {
			prefix := groupPrefix(k)
			if lastPrefix != "" && prefix != lastPrefix {
				if _, err := fmt.Fprintln(w); err != nil {
					return err
				}
			}
			lastPrefix = prefix
		}
		v := secrets[k]
		if needsQuotes(v) {
			v = `"` + strings.ReplaceAll(v, `"`, `\"`) + `"`
		}
		if _, err := fmt.Fprintf(w, "%s=%s\n", k, v); err != nil {
			return err
		}
	}
	return nil
}

func groupPrefix(key string) string {
	if idx := strings.Index(key, "_"); idx > 0 {
		return key[:idx]
	}
	return key
}

func needsQuotes(v string) bool {
	return strings.ContainsAny(v, " \t\n#")
}
