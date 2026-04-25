package envwildcard_test

import (
	"fmt"
	"sort"

	"github.com/vaultline/vaultline/internal/envwildcard"
)

func ExampleMatcher_Apply() {
	secrets := map[string]string{
		"DB_HOST":     "localhost",
		"DB_PASSWORD": "s3cr3t",
		"APP_NAME":    "vaultline",
		"APP_SECRET":  "topsecret",
	}

	m, _ := envwildcard.New(
		[]string{"DB_*", "APP_*"},
		[]string{"*_PASSWORD", "*_SECRET"},
	)

	out, _ := m.Apply(secrets)

	keys := make([]string, 0, len(out))
	for k := range out {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Printf("%s=%s\n", k, out[k])
	}
	// Output:
	// APP_NAME=vaultline
	// DB_HOST=localhost
}
