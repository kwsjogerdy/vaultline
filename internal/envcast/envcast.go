// Package envcast provides type-casting utilities for secret values.
package envcast

import (
	"fmt"
	"strconv"
	"strings"
)

// Caster converts secret string values to typed Go values.
type Caster struct{}

// New returns a new Caster.
func New() *Caster { return &Caster{} }

// AsString returns the value as-is.
func (c *Caster) AsString(secrets map[string]string, key string) (string, error) {
	v, ok := secrets[key]
	if !ok {
		return "", fmt.Errorf("envcast: key %q not found", key)
	}
	return v, nil
}

// AsInt parses the value as an integer.
func (c *Caster) AsInt(secrets map[string]string, key string) (int64, error) {
	v, err := c.AsString(secrets, key)
	if err != nil {
		return 0, err
	}
	n, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("envcast: key %q cannot be cast to int: %w", key, err)
	}
	return n, nil
}

// AsBool parses the value as a boolean.
func (c *Caster) AsBool(secrets map[string]string, key string) (bool, error) {
	v, err := c.AsString(secrets, key)
	if err != nil {
		return false, err
	}
	b, err := strconv.ParseBool(strings.TrimSpace(v))
	if err != nil {
		return false, fmt.Errorf("envcast: key %q cannot be cast to bool: %w", key, err)
	}
	return b, nil
}

// AsFloat parses the value as a float64.
func (c *Caster) AsFloat(secrets map[string]string, key string) (float64, error) {
	v, err := c.AsString(secrets, key)
	if err != nil {
		return 0, err
	}
	f, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
	if err != nil {
		return 0, fmt.Errorf("envcast: key %q cannot be cast to float: %w", key, err)
	}
	return f, nil
}

// CastAll attempts to cast all values using the provided type map (key -> "int"|"bool"|"float"|"string").
// Returns a map of key -> cast value as interface{}, and a slice of errors.
func (c *Caster) CastAll(secrets map[string]string, types map[string]string) (map[string]interface{}, []error) {
	out := make(map[string]interface{}, len(secrets))
	var errs []error
	for k, v := range secrets {
		t, ok := types[k]
		if !ok {
			out[k] = v
			continue
		}
		switch t {
		case "int":
			n, err := c.AsInt(secrets, k)
			if err != nil {
				errs = append(errs, err)
			} else {
				out[k] = n
			}
		case "bool":
			b, err := c.AsBool(secrets, k)
			if err != nil {
				errs = append(errs, err)
			} else {
				out[k] = b
			}
		case "float":
			f, err := c.AsFloat(secrets, k)
			if err != nil {
				errs = append(errs, err)
			} else {
				out[k] = f
			}
		default:
			out[k] = v
		}
	}
	return out, errs
}
