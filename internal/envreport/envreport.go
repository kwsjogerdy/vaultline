// Package envreport generates summary reports of secret sets.
package envreport

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// Format defines the output format for a report.
type Format string

const (
	FormatText Format = "text"
	FormatMarkdown Format = "markdown"
)

// Report holds statistics about a set of secrets.
type Report struct {
	Total    int
	Keys     []string
	BySuffix map[string]int
}

// Reporter generates env reports.
type Reporter struct {
	w      io.Writer
	format Format
}

// New creates a new Reporter writing to w.
func New(w io.Writer, format Format) *Reporter {
	return &Reporter{w: w, format: format}
}

// Build computes a Report from a secrets map.
func Build(secrets map[string]string) Report {
	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	bySuffix := map[string]int{}
	for _, k := range keys {
		parts := strings.SplitN(k, "_", 2)
		if len(parts) > 1 {
			bySuffix[parts[0]]++
		} else {
			bySuffix["(none)"]++
		}
	}
	return Report{Total: len(keys), Keys: keys, BySuffix: bySuffix}
}

// Write renders the report to the reporter's writer.
func (r *Reporter) Write(rep Report) error {
	switch r.format {
	case FormatMarkdown:
		return r.writeMarkdown(rep)
	default:
		return r.writeText(rep)
	}
}

func (r *Reporter) writeText(rep Report) error {
	fmt.Fprintf(r.w, "Total secrets: %d\n", rep.Total)
	fmt.Fprintf(r.w, "Keys:\n")
	for _, k := range rep.Keys {
		fmt.Fprintf(r.w, "  - %s\n", k)
	}
	fmt.Fprintf(r.w, "By prefix:\n")
	for _, prefix := range sortedKeys(rep.BySuffix) {
		fmt.Fprintf(r.w, "  %s: %d\n", prefix, rep.BySuffix[prefix])
	}
	return nil
}

func (r *Reporter) writeMarkdown(rep Report) error {
	fmt.Fprintf(r.w, "## Secret Report\n\n")
	fmt.Fprintf(r.w, "**Total:** %d\n\n", rep.Total)
	fmt.Fprintf(r.w, "| Key |\n|-----|\n")
	for _, k := range rep.Keys {
		fmt.Fprintf(r.w, "| %s |\n", k)
	}
	fmt.Fprintf(r.w, "\n**By Prefix:**\n\n")
	for _, prefix := range sortedKeys(rep.BySuffix) {
		fmt.Fprintf(r.w, "- %s: %d\n", prefix, rep.BySuffix[prefix])
	}
	return nil
}

func sortedKeys(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
