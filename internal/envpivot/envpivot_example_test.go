package envpivot_test

import (
	"fmt"
	"sort"

	"github.com/vaultline/vaultline/internal/envpivot"
)

func ExamplePivoter_Apply() {
	p, _ := envpivot.New(envpivot.Options{Collision: envpivot.CollisionError})
	secrets := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	}
	result, _ := p.Apply(secrets)

	keys := make([]string, 0, len(result))
	for k := range result {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf("%s -> %s\n", k, result[k])
	}
	// Output:
	// 5432 -> DB_PORT
	// localhost -> DB_HOST
}
