// Package envpatch applies partial updates to a secret map using a sequence
// of typed operations: set, delete, and rename. It is safe to use with any
// map[string]string secret store and does not mutate the input.
package envpatch
