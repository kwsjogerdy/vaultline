package envgraph_test

import (
	"fmt"

	"github.com/vaultline/vaultline/internal/envgraph"
)

func ExampleGraph_DOT() {
	secrets := map[string]string{
		"HOST":    "db.example.com",
		"DB_URL":  "postgres://${HOST}/mydb",
		"APP_DSN": "${DB_URL}?sslmode=disable",
	}
	g := envgraph.New(secrets)
	// Print roots (keys nothing else references)
	fmt.Println("roots:", g.Roots())
	// Output:
	// roots: [APP_DSN]
}
