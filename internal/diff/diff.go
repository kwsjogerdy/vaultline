package diff

import "fmt"

// Result holds the comparison between existing and incoming secrets.
type Result struct {
	Added     map[string]string
	Removed   map[string]string
	Changed   map[string]string
	Unchanged map[string]string
}

// Compare compares existing env vars against incoming secrets and returns a Result.
func Compare(existing, incoming map[string]string) Result {
	r := Result{
		Added:     make(map[string]string),
		Removed:   make(map[string]string),
		Changed:   make(map[string]string),
		Unchanged: make(map[string]string),
	}

	for k, v := range incoming {
		oldVal, ok := existing[k]
		if !ok {
			r.Added[k] = v
		} else if oldVal != v {
			r.Changed[k] = v
		} else {
			r.Unchanged[k] = v
		}
	}

	for k, v := range existing {
		if _, ok := incoming[k]; !ok {
			r.Removed[k] = v
		}
	}

	return r
}

// HasChanges returns true if there are any added, removed, or changed keys.
func (r Result) HasChanges() bool {
	return len(r.Added) > 0 || len(r.Removed) > 0 || len(r.Changed) > 0
}

// Summary returns a human-readable summary string.
func (r Result) Summary() string {
	return fmt.Sprintf("+%d added, ~%d changed, -%d removed, %d unchanged",
		len(r.Added), len(r.Changed), len(r.Removed), len(r.Unchanged))
}
