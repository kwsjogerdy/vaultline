package filter

import "strings"

// Rule defines inclusion/exclusion rules for secret keys.
type Rule struct {
	Prefixes []string
	Excludes []string
}

// Filter applies rules to a map of secrets and returns only the allowed keys.
type Filter struct {
	rule Rule
}

// New creates a new Filter with the given Rule.
func New(rule Rule) *Filter {
	return &Filter{rule: rule}
}

// Apply returns a filtered copy of secrets based on the filter rules.
// If no prefixes are defined, all keys are included unless excluded.
func (f *Filter) Apply(secrets map[string]string) map[string]string {
	result := make(map[string]string)

	for k, v := range secrets {
		if f.isExcluded(k) {
			continue
		}
		if len(f.rule.Prefixes) == 0 || f.hasPrefix(k) {
			result[k] = v
		}
	}

	return result
}

func (f *Filter) isExcluded(key string) bool {
	for _, ex := range f.rule.Excludes {
		if strings.EqualFold(key, ex) {
			return true
		}
	}
	return false
}

func (f *Filter) hasPrefix(key string) bool {
	for _, p := range f.rule.Prefixes {
		if strings.HasPrefix(strings.ToUpper(key), strings.ToUpper(p)) {
			return true
		}
	}
	return false
}
