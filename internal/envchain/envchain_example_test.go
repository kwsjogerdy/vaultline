package envchain_test

import (
	"fmt"

	"github.com/vaultline/vaultline/internal/envchain"
)

func ExampleChain_Resolve() {
	chain := envchain.New([]envchain.Source{
		{Name: "defaults", Priority: 1, Secrets: map[string]string{"LOG_LEVEL": "info", "PORT": "8080"}},
		{Name: "vault", Priority: 5, Secrets: map[string]string{"PORT": "9090", "DB_PASS": "s3cr3t"}},
	})
	resolved := chain.Resolve()
	fmt.Println(resolved["PORT"])    // vault wins
	fmt.Println(resolved["LOG_LEVEL"]) // defaults only
	// Output:
	// 9090
	// info
}
