// Package envdiff provides side-by-side comparison of two env sources.
package envdiff

import "fmt"

// Entry represents a single compared key across two sources.
type Entry struct {
	Key    string
	Left   string
	Right  string
	Status string // "match", "changed", "left_only", "right_only"
}

// Result holds the full comparison output.
type Result struct {
	Entries []Entry
}

// HasDrift returns true if any entry is not a match.
func (r Result) HasDrift() bool {
	for _, e := range r.Entries {
		if e.Status != "match" {
			return true
		}
	}
	return false
}

// Summary returns a human-readable summary string.
func (r Result) Summary() string {
	var changed, leftOnly, rightOnly, match int
	for _, e := range r.Entries {
		switch e.Status {
		case "match":
			match++
		case "changed":
			changed++
		case "left_only":
			leftOnly++
		case "right_only":
			rightOnly++
		}
	}
	return fmt.Sprintf("match=%d changed=%d left_only=%d right_only=%d", match, changed, leftOnly, rightOnly)
}

// Compare performs a key-by-key comparison of two secret maps.
func Compare(left, right map[string]string) Result {
	seen := map[string]bool{}
	var entries []Entry

	for k, lv := range left {
		seen[k] = true
		if rv, ok := right[k]; ok {
			status := "match"
			if lv != rv {
				status = "changed"
			}
			entries = append(entries, Entry{Key: k, Left: lv, Right: rv, Status: status})
		} else {
			entries = append(entries, Entry{Key: k, Left: lv, Status: "left_only"})
		}
	}

	for k, rv := range right {
		if !seen[k] {
			entries = append(entries, Entry{Key: k, Right: rv, Status: "right_only"})
		}
	}

	return Result{Entries: entries}
}
