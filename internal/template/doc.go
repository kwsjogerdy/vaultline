// Package template provides Go text/template-based rendering for .env files.
// Secrets fetched from Vault are passed as a map[string]string to the template,
// accessible via `index . "KEY"`. Helper functions `default` and `required`
// allow safe handling of optional and mandatory secret values.
package template
