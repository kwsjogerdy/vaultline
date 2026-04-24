// Package envpatch applies partial updates (patches) to a secret map,
// supporting set, delete, and rename operations in a single pass.
package envpatch

import (
	"errors"
	"fmt"	
)

// OpType represents the kind of patch operation.
type OpType string

const (
	OpSet    OpType = "set"
	OpDelete OpType = "delete"
	OpRename OpType = "rename"
)

// Op describes a single patch operation.
type Op struct {
	Type  OpType
	Key   string
	Value string // used by set
	NewKey string // used by rename
}

// Result holds the outcome of applying a patch.
type Result struct {
	Applied []string
	Skipped []string
}

// Patcher applies a sequence of Ops to a secret map.
type Patcher struct {
	ops []Op
}

// New creates a Patcher with the given operations.
func New(ops []Op) (*Patcher, error) {
	for i, op := range ops {
		if op.Key == "" {
			return nil, fmt.Errorf("op[%d]: key must not be empty", i)
		}
		if op.Type == OpRename && op.NewKey == "" {
			return nil, fmt.Errorf("op[%d]: rename requires a non-empty new_key", i)
		}
		if op.Type != OpSet && op.Type != OpDelete && op.Type != OpRename {
			return nil, fmt.Errorf("op[%d]: unknown op type %q", i, op.Type)
		}
	}
	return &Patcher{ops: ops}, nil
}

// Apply executes the patch operations against secrets and returns a new map.
func (p *Patcher) Apply(secrets map[string]string) (map[string]string, Result, error) {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}

	var res Result
	for _, op := range p.ops {
		switch op.Type {
		case OpSet:
			out[op.Key] = op.Value
			res.Applied = append(res.Applied, op.Key)
		case OpDelete:
			if _, ok := out[op.Key]; !ok {
				res.Skipped = append(res.Skipped, op.Key)
				continue
			}
			delete(out, op.Key)
			res.Applied = append(res.Applied, op.Key)
		case OpRename:
			val, ok := out[op.Key]
			if !ok {
				res.Skipped = append(res.Skipped, op.Key)
				continue
			}
			out[op.NewKey] = val
			delete(out, op.Key)
			res.Applied = append(res.Applied, op.Key)
		default:
			return nil, res, errors.New("unexpected op type")
		}
	}
	return out, res, nil
}
