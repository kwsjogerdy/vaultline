package envsanitize_test

import (
	"fmt"

	"github.com/vaultline/vaultline/internal/envsanitize"
)

func ExampleSanitizer_Apply() {
	s, err := envsanitize.New(envsanitize.Options{
		StripNewlines: true,
		TrimSpace:     true,
		Rules: []envsanitize.Rule{
			{Pattern: `\s{2,}`, Replacement: " "},
		},
	})
	if err != nil {
		panic(err)
	}

	secrets := map[string]string{
		"DB_PASS": "  my secret\nvalue  ",
		"API_KEY": "abc123",
	}

	out := s.Apply(secrets)
	fmt.Println(out["DB_PASS"])
	fmt.Println(out["API_KEY"])
	// Output:
	// my secret value
	// abc123
}
