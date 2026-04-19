// Package envfence restricts which keys are allowed to be written
// based on an allowlist or denylist policy.
package envfence

import "fmt"

// Policy controls whether the fence operates as an allowlist or denylist.
type Policy string

const (
	PolicyAllow Policy = "allow"
	PolicyDeny  Policy = "deny"
)

// Fence enforces key access policies on a secret map.
type Fence struct {
	policy Policy
	keys   map[string]struct{}
}

// New creates a Fence with the given policy and key list.
func New(policy Policy, keys []string) (*Fence, error) {
	if policy != PolicyAllow && policy != PolicyDeny {
		return nil, fmt.Errorf("envfence: unknown policy %q", policy)
	}
	km := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		km[k] = struct{}{}
	}
	return &Fence{policy: policy, keys: km}, nil
}

// Apply returns a filtered copy of secrets according to the fence policy.
func (f *Fence) Apply(secrets map[string]string) map[string]string {
	out := make(map[string]string)
	for k, v := range secrets {
		_, listed := f.keys[k]
		switch f.policy {
		case PolicyAllow:
			if listed {
				out[k] = v
			}
		case PolicyDeny:
			if !listed {
				out[k] = v
			}
		}
	}
	return out
}

// Blocked returns keys present in secrets that are blocked by the fence.
func (f *Fence) Blocked(secrets map[string]string) []string {
	var blocked []string
	for k := range secrets {
		_, listed := f.keys[k]
		if f.policy == PolicyAllow && !listed {
			blocked = append(blocked, k)
		} else if f.policy == PolicyDeny && listed {
			blocked = append(blocked, k)
		}
	}
	return blocked
}
