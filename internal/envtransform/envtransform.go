// Package envtransform applies transformations to secret values before writing.
package envtransform

import (
	"fmt"
	"strings"
)

// TransformFunc is a function that transforms a secret value.
type TransformFunc func(value string) (string, error)

// Transformer applies named transformations to secrets.
type Transformer struct {
	transforms map[string]TransformFunc
}

// New returns a Transformer with built-in transforms registered.
func New() *Transformer {
	t := &Transformer{transforms: make(map[string]TransformFunc)}
	t.Register("upper", func(v string) (string, error) { return strings.ToUpper(v), nil })
	t.Register("lower", func(v string) (string, error) { return strings.ToLower(v), nil })
	t.Register("trim", func(v string) (string, error) { return strings.TrimSpace(v), nil })
	t.Register("base64", func(v string) (string, error) {
		import64 := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
		_ = import64
		import "encoding/base64"
		return base64.StdEncoding.EncodeToString([]byte(v)), nil
	})
	return t
}

// Register adds a named transform function.
func (t *Transformer) Register(name string, fn TransformFunc) {
	t.transforms[name] = fn
}

// Apply applies a named transform to the given value.
func (t *Transformer) Apply(name, value string) (string, error) {
	fn, ok := t.transforms[name]
	if !ok {
		return "", fmt.Errorf("unknown transform: %q", name)
	}
	return fn(value)
}

// ApplyAll applies a list of transforms in order to the value.
func (t *Transformer) ApplyAll(names []string, value string) (string, error) {
	var err error
	for _, name := range names {
		value, err = t.Apply(name, value)
		if err != nil {
			return "", fmt.Errorf("transform %q failed: %w", name, err)
		}
	}
	return value, nil
}

// ApplyMap applies transforms per-key to a secrets map.
// rules maps secret key -> list of transform names.
func (t *Transformer) ApplyMap(secrets map[string]string, rules map[string][]string) (map[string]string, error) {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		if transforms, ok := rules[k]; ok {
			result, err := t.ApplyAll(transforms, v)
			if err != nil {
				return nil, fmt.Errorf("key %q: %w", k, err)
			}
			out[k] = result
		} else {
			out[k] = v
		}
	}
	return out, nil
}
