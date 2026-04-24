package envindex_test

import (
	"fmt"

	"github.com/vaultline/vaultline/internal/envindex"
)

func ExampleIndex_Search() {
	secrets := map[string]string{
		"DATABASE_URL": "postgres://localhost/mydb",
		"API_KEY":      "s3cr3t",
		"REDIS_HOST":   "localhost",
	}
	labels := map[string]string{
		"API_KEY": "auth",
	}

	idx := envindex.New(secrets, labels)
	results := idx.Search("auth")
	for _, r := range results {
		fmt.Printf("%s matched on %s\n", r.Key, r.MatchedOn)
	}
	// Output:
	// API_KEY matched on label
}
