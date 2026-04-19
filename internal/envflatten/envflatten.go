// Package envflatten flattens nested map structures into dot-notation env keys.
package envflatten

import (
	"fmt"
	"strings"
)

// Flattener converts nested maps to flat key=value pairs.
type Flattener struct {
	separator string
	prefix    string
}

// New returns a Flattener with the given separator (e.g. "__" or ".").
func New(separator string) *Flattener {
	if separator == "" {
		separator = "__"
	}
	return &Flattener{separator: separator}
}

// WithPrefix returns a new Flattener that prepends a prefix to all keys.
func (f *Flattener) WithPrefix(prefix string) *Flattener {
	return &Flattener{separator: f.separator, prefix: prefix}
}

// Flatten converts a nested map[string]any into a flat map[string]string.
func (f *Flattener) Flatten(input map[string]any) (map[string]string, error) {
	result := make(map[string]string)
	if err := flatten(input, f.prefix, f.separator, result); err != nil {
		return nil, err
	}
	return result, nil
}

func flatten(input map[string]any, prefix, sep string, out map[string]string) error {
	for k, v := range input {
		key := k
		if prefix != "" {
			key = prefix + sep + k
		}
		key = strings.ToUpper(key)
		switch val := v.(type) {
		case map[string]any:
			if err := flatten(val, key, sep, out); err != nil {
				return err
			}
		case string:
			out[key] = val
		case int, int64, float64, bool:
			out[key] = fmt.Sprintf("%v", val)
		case nil:
			out[key] = ""
		default:
			return fmt.Errorf("unsupported value type for key %q: %T", key, v)
		}
	}
	return nil
}
