// Package envtemplate renders Go text/template strings using secrets as the data source.
package envtemplate

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"text/template"
)

// Renderer applies Go templates to a secrets map.
type Renderer struct {
	secrets map[string]string
}

// New creates a Renderer backed by the given secrets map.
func New(secrets map[string]string) *Renderer {
	copy := make(map[string]string, len(secrets))
	for k, v := range secrets {
		copy[k] = v
	}
	return &Renderer{secrets: copy}
}

// Render executes the given template text with secrets as the data context.
// Template variables are accessed via {{.KEY}} or the helper func `secret "KEY"`.
func (r *Renderer) Render(text string) (string, error) {
	funcMap := template.FuncMap{
		"secret": func(key string) (string, error) {
			v, ok := r.secrets[key]
			if !ok {
				return "", fmt.Errorf("secret %q not found", key)
			}
			return v, nil
		},
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
	}

	tmpl, err := template.New("env").Option("missingkey=error").Funcs(funcMap).Parse(text)
	if err != nil {
		return "", fmt.Errorf("parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, r.secrets); err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}
	return buf.String(), nil
}

// RenderMap applies Render to every value in the secrets map and returns a new map.
func (r *Renderer) RenderMap() (map[string]string, error) {
	out := make(map[string]string, len(r.secrets))
	keys := make([]string, 0, len(r.secrets))
	for k := range r.secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		rendered, err := r.Render(r.secrets[k])
		if err != nil {
			return nil, fmt.Errorf("key %q: %w", k, err)
		}
		out[k] = rendered
	}
	return out, nil
}
