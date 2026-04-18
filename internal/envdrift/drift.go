// Package envdrift detects drift between local .env files and Vault secrets.
package envdrift

import "fmt"

// Status represents the drift state of a single key.
type Status int

const (
	StatusMatch   Status = iota // local value matches Vault
	StatusDrifted               // local value differs from Vault
	StatusMissing               // key missing locally
	StatusExtra                 // key present locally but not in Vault
)

func (s Status) String() string {
	switch s {
	case StatusMatch:
		return "match"
	case StatusDrifted:
		return "drifted"
	case StatusMissing:
		return "missing"
	case StatusExtra:
		return "extra"
	default:
		return "unknown"
	}
}

// Entry describes the drift status of one key.
type Entry struct {
	Key    string
	Status Status
}

// Report holds all drift entries for a comparison.
type Report struct {
	Entries []Entry
}

// HasDrift returns true if any key is not StatusMatch.
func (r *Report) HasDrift() bool {
	for _, e := range r.Entries {
		if e.Status != StatusMatch {
			return true
		}
	}
	return false
}

// Summary returns a human-readable drift summary.
func (r *Report) Summary() string {
	counts := map[Status]int{}
	for _, e := range r.Entries {
		counts[e.Status]++
	}
	return fmt.Sprintf("match=%d drifted=%d missing=%d extra=%d",
		counts[StatusMatch], counts[StatusDrifted], counts[StatusMissing], counts[StatusExtra])
}

// Detect compares local env map against vault secrets map and returns a Report.
func Detect(local, vault map[string]string) *Report {
	report := &Report{}
	for k, vaultVal := range vault {
		if localVal, ok := local[k]; !ok {
			report.Entries = append(report.Entries, Entry{Key: k, Status: StatusMissing})
		} else if localVal != vaultVal {
			report.Entries = append(report.Entries, Entry{Key: k, Status: StatusDrifted})
		} else {
			report.Entries = append(report.Entries, Entry{Key: k, Status: StatusMatch})
		}
	}
	for k := range local {
		if _, ok := vault[k]; !ok {
			report.Entries = append(report.Entries, Entry{Key: k, Status: StatusExtra})
		}
	}
	return report
}
