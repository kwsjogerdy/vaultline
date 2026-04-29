// Package envsanitize provides a configurable sanitizer for cleaning secret
// values before writing them to .env files or passing them to downstream tools.
// It supports stripping null bytes, collapsing newlines, trimming whitespace,
// and applying arbitrary regexp-based replacement rules.
package envsanitize
