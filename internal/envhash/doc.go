// Package envhash provides stable SHA-256 fingerprinting for secret maps.
// It is used to detect whether secrets have changed between sync operations
// without persisting or comparing raw secret values directly.
package envhash
