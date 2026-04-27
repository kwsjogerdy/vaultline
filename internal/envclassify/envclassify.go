// Package envclassify categorises secret keys into semantic buckets
// such as credentials, endpoints, feature flags, and identifiers.
package envclassify

import (
	"sort"
	"strings"
)

// Category represents a semantic classification label.
type Category string

const (
	CategoryCredential Category = "credential"
	CategoryEndpoint   Category = "endpoint"
	CategoryFeatureFlag Category = "feature_flag"
	CategoryIdentifier Category = "identifier"
	CategoryUnknown    Category = "unknown"
)

// Result holds the classification for a single key.
type Result struct {
	Key      string
	Value    string
	Category Category
}

// Classifier assigns categories to secret keys.
type Classifier struct {
	rules []rule
}

type rule struct {
	category Category
	patterns []string
}

var defaultRules = []rule{
	{CategoryCredential, []string{"password", "passwd", "secret", "token", "api_key", "apikey", "auth", "credential", "private_key", "passphrase"}},
	{CategoryEndpoint, []string{"url", "host", "endpoint", "addr", "address", "port", "dsn", "uri", "base_url"}},
	{CategoryFeatureFlag, []string{"enable_", "disable_", "feature_", "flag_", "_enabled", "_disabled"}},
	{CategoryIdentifier, []string{"id", "_id", "uuid", "name", "project", "region", "zone", "namespace", "env"}},
}

// New returns a Classifier using the default rule set.
func New() *Classifier {
	return &Classifier{rules: defaultRules}
}

// Classify returns a Result slice for the provided secrets map.
func (c *Classifier) Classify(secrets map[string]string) []Result {
	results := make([]Result, 0, len(secrets))
	for k, v := range secrets {
		results = append(results, Result{Key: k, Value: v, Category: c.categorise(k)})
	}
	sort.Slice(results, func(i, j int) bool { return results[i].Key < results[j].Key })
	return results
}

// ByCategory groups results into a map keyed by Category.
func ByCategory(results []Result) map[Category][]Result {
	out := map[Category][]Result{}
	for _, r := range results {
		out[r.Category] = append(out[r.Category], r)
	}
	return out
}

func (c *Classifier) categorise(key string) Category {
	lower := strings.ToLower(key)
	for _, r := range c.rules {
		for _, p := range r.patterns {
			if strings.Contains(lower, p) {
				return r.category
			}
		}
	}
	return CategoryUnknown
}
