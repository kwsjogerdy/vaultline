// Package envjoin merges two secret maps using a join key, producing a
// combined map where values from both sides are concatenated or composed.
package envjoin

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

// JoinMode controls how values from left and right maps are combined.
type JoinMode string

const (
	// ModeConcat joins values with a separator string.
	ModeConcat JoinMode = "concat"
	// ModeLeft keeps only the left value when a key exists in both maps.
	ModeLeft JoinMode = "left"
	// ModeRight keeps only the right value when a key exists in both maps.
	ModeRight JoinMode = "right"
)

// Joiner combines two secret maps.
type Joiner struct {
	mode JoinMode
	sep  string
}

// New creates a Joiner with the given mode and separator.
// separator is only used for ModeConcat.
func New(mode JoinMode, sep string) (*Joiner, error) {
	switch mode {
	case ModeConcat, ModeLeft, ModeRight:
	default:
		return nil, fmt.Errorf("envjoin: unknown join mode %q", mode)
	}
	return &Joiner{mode: mode, sep: sep}, nil
}

// Join merges left and right into a single map.
// Keys present in only one side are always included as-is.
func (j *Joiner) Join(left, right map[string]string) (map[string]string, error) {
	if left == nil {
		return nil, errors.New("envjoin: left map must not be nil")
	}
	if right == nil {
		return nil, errors.New("envjoin: right map must not be nil")
	}

	out := make(map[string]string, len(left))
	for k, v := range left {
		out[k] = v
	}

	for k, rv := range right {
		lv, exists := out[k]
		if !exists {
			out[k] = rv
			continue
		}
		switch j.mode {
		case ModeConcat:
			out[k] = lv + j.sep + rv
		case ModeLeft:
			// keep lv as-is
		case ModeRight:
			out[k] = rv
		}
	}
	return out, nil
}

// Keys returns a sorted slice of all keys present in the joined result.
func Keys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Summary returns a human-readable summary of the join result.
func Summary(left, right, result map[string]string) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "left=%d right=%d result=%d\n", len(left), len(right), len(result))
	return sb.String()
}
