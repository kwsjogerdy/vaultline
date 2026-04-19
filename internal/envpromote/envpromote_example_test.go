package envpromote_test

import (
	"fmt"

	"vaultline/internal/envpromote"
)

func ExamplePromoter_Promote() {
	src := map[string]string{
		"DEV_DB_HOST": "localhost",
		"DEV_API_KEY": "dev-secret",
	}
	dst := map[string]string{}

	p := envpromote.New("dev", "staging").WithPrefix("DEV_")
	result, info := p.Promote(src, dst, nil)

	fmt.Println(result["DB_HOST"])
	fmt.Println(len(info.Copied))
	// Output:
	// localhost
	// 2
}
