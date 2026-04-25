package envstats_test

import (
	"fmt"

	"github.com/vaultline/vaultline/internal/envstats"
)

func ExampleAnalyzer_Compute() {
	a := envstats.New("_")
	secrets := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"APP_KEY": "s3cr3t",
	}
	s := a.Compute(secrets)
	fmt.Println(s.Total)
	fmt.Println(s.PrefixCounts["DB"])
	// Output:
	// 3
	// 2
}
