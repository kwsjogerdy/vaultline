// Package envpivot transposes a secrets map by swapping keys and values.
// Collision strategies (error, skip, suffix) control behaviour when
// multiple keys share the same value.
package envpivot
