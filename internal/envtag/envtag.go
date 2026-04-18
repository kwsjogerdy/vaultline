// Package envtag provides tagging and grouping of secrets by environment label.
package envtag

import (
	"fmt"
	"strings"
)

// Tag represents an environment label attached to a set of secret keys.
type Tag struct {
	Label string
	Keys  []string
}

// Tagger groups secrets by their tag prefix (e.g. "prod/", "dev/").
type Tagger struct {
	tags map[string]*Tag
}

// New creates a new Tagger.
func New() *Tagger {
	return &Tagger{tags: make(map[string]*Tag)}
}

// Add registers a key under the given label.
func (t *Tagger) Add(label, key string) {
	label = strings.ToLower(strings.TrimSpace(label))
	if label == "" {
		label = "default"
	}
	if _, ok := t.tags[label]; !ok {
		t.tags[label] = &Tag{Label: label}
	}
	t.tags[label].Keys = append(t.tags[label].Keys, key)
}

// InferAndAdd inspects the key for a known prefix (e.g. "PROD_", "DEV_") and tags accordingly.
func (t *Tagger) InferAndAdd(key string) {
	upper := strings.ToUpper(key)
	for _, env := range []string{"PROD", "STAGING", "DEV", "TEST"} {
		if strings.HasPrefix(upper, env+"_") {
			t.Add(strings.ToLower(env), key)
			return
		}
	}
	t.Add("default", key)
}

// Get returns the Tag for the given label, or an error if not found.
func (t *Tagger) Get(label string) (*Tag, error) {
	label = strings.ToLower(strings.TrimSpace(label))
	tag, ok := t.tags[label]
	if !ok {
		return nil, fmt.Errorf("envtag: label %q not found", label)
	}
	return tag, nil
}

// All returns all registered tags.
func (t *Tagger) All() []*Tag {
	out := make([]*Tag, 0, len(t.tags))
	for _, tag := range t.tags {
		out = append(out, tag)
	}
	return out
}

// Filter returns only the secrets whose keys belong to the given label.
func (t *Tagger) Filter(label string, secrets map[string]string) (map[string]string, error) {
	tag, err := t.Get(label)
	if err != nil {
		return nil, err
	}
	result := make(map[string]string)
	for _, k := range tag.Keys {
		if v, ok := secrets[k]; ok {
			result[k] = v
		}
	}
	return result, nil
}
