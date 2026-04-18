// Package envresolve resolves final secret values by interpolating
// references between keys within a secret map.
package envresolve

import (
	"fmt"
	"regexp"
	"strings"
)

var refPattern = regexp.MustCompile(`\$\{([A-Z0-9_]+)\}`)

// Resolver expands ${OTHER_KEY} references inside secret values.
type Resolver struct {
	maxDepth int
}

// New returns a Resolver with a default max interpolation depth of 10.
func New() *Resolver {
	return &Resolver{maxDepth: 10}
}

// Resolve returns a new map where all ${KEY} references are replaced
// with the value of that key from the same map. Circular or missing
// references return an error.
func (r *Resolver) Resolve(secrets map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		resolved, err := r.expand(v, secrets, 0)
		if err != nil {
			return nil, fmt.Errorf("key %q: %w", k, err)
		}
		out[k] = resolved
	}
	return out, nil
}

func (r *Resolver) expand(val string, secrets map[string]string, depth int) (string, error) {
	if depth > r.maxDepth {
		return "", fmt.Errorf("max interpolation depth (%d) exceeded", r.maxDepth)
	}
	if !strings.Contains(val, "${") {
		return val, nil
	}
	var expandErr error
	result := refPattern.ReplaceAllStringFunc(val, func(match string) string {
		if expandErr != nil {
			return ""
		}
		key := refPattern.FindStringSubmatch(match)[1]
		replacement, ok := secrets[key]
		if !ok {
			expandErr = fmt.Errorf("undefined reference ${%s}", key)
			return ""
		}
		expanded, err := r.expand(replacement, secrets, depth+1)
		if err != nil {
			expandErr = err
			return ""
		}
		return expanded
	})
	if expandErr != nil {
		return "", expandErr
	}
	return result, nil
}
