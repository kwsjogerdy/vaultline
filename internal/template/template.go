// Package template renders .env files from Go text/template strings,
// allowing dynamic secret injection with optional default values.
package template

import (
	"bytes"
	"fmt"
	"os"
	"text/template"
)

// Renderer renders a template file using secrets as the data source.
type Renderer struct {
	templatePath string
}

// New returns a new Renderer for the given template file path.
func New(templatePath string) *Renderer {
	return &Renderer{templatePath: templatePath}
}

// Render reads the template file, executes it with the provided secrets map,
// and returns the rendered output as a string.
func (r *Renderer) Render(secrets map[string]string) (string, error) {
	raw, err := os.ReadFile(r.templatePath)
	if err != nil {
		return "", fmt.Errorf("template: read %q: %w", r.templatePath, err)
	}
	return RenderBytes(raw, secrets)
}

// RenderBytes executes a template from raw bytes with the provided secrets map.
func RenderBytes(raw []byte, secrets map[string]string) (string, error) {
	funcMap := template.FuncMap{
		"default": func(def, val string) string {
			if val == "" {
				return def
			}
			return val
		},
		"required": func(key, val string) (string, error) {
			if val == "" {
				return "", fmt.Errorf("template: required key %q is missing or empty", key)
			}
			return val, nil
		},
	}

	tmpl, err := template.New("env").Funcs(funcMap).Parse(string(raw))
	if err != nil {
		return "", fmt.Errorf("template: parse: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, secrets); err != nil {
		return "", fmt.Errorf("template: execute: %w", err)
	}
	return buf.String(), nil
}
