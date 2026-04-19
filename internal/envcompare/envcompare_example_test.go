package envcompare_test

import (
	"fmt"

	"github.com/vaultline/vaultline/internal/envcompare"
)

func ExampleComparer_Compare() {
	c := envcompare.New("dev", "prod")
	report := c.Compare(
		map[string]string{"DB_HOST": "localhost", "API_KEY": "abc"},
		map[string]string{"DB_HOST": "prod.db", "API_KEY": "abc"},
	)
	fmt.Println(report.Summary())
	// Output: match=1 differ=1 left_only=0 right_only=0
}
