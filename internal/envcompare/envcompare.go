// Package envcompare compares secrets across two Vault paths or environments.
package envcompare

import "fmt"

// Result holds the comparison outcome for a single key.
type Result struct {
	Key    string
	Left   string
	Right  string
	Status string // "match", "differ", "left_only", "right_only"
}

// Report is the full comparison output.
type Report struct {
	Results []Result
	Left    string
	Right   string
}

// Comparer compares two secret maps.
type Comparer struct {
	leftLabel  string
	rightLabel string
}

// New creates a Comparer with labels for each side.
func New(leftLabel, rightLabel string) *Comparer {
	return &Comparer{leftLabel: leftLabel, rightLabel: rightLabel}
}

// Compare returns a Report for two secret maps.
func (c *Comparer) Compare(left, right map[string]string) Report {
	seen := map[string]bool{}
	var results []Result

	for k, lv := range left {
		seen[k] = true
		if rv, ok := right[k]; ok {
			status := "match"
			if lv != rv {
				status = "differ"
			}
			results = append(results, Result{Key: k, Left: lv, Right: rv, Status: status})
		} else {
			results = append(results, Result{Key: k, Left: lv, Status: "left_only"})
		}
	}

	for k, rv := range right {
		if !seen[k] {
			results = append(results, Result{Key: k, Right: rv, Status: "right_only"})
		}
	}

	return Report{Results: results, Left: c.leftLabel, Right: c.rightLabel}
}

// Summary returns a human-readable summary string.
func (r Report) Summary() string {
	var match, differ, leftOnly, rightOnly int
	for _, res := range r.Results {
		switch res.Status {
		case "match":
			match++
		case "differ":
			differ++
		case "left_only":
			leftOnly++
		case "right_only":
			rightOnly++
		}
	}
	return fmt.Sprintf("match=%d differ=%d left_only=%d right_only=%d", match, differ, leftOnly, rightOnly)
}
