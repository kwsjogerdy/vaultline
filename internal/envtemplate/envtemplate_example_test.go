package envtemplate_test

import (
	"fmt"

	"github.com/vaultline/vaultline/internal/envtemplate"
)

func ExampleRenderer_Render() {
	secrets := map[string]string{
		"DB_HOST": "db.internal",
		"DB_PORT": "5432",
	}
	r := envtemplate.New(secrets)
	out, err := r.Render(`host={{ secret "DB_HOST" }} port={{ secret "DB_PORT" }}`)
	if err != nil {
		panic(err)
	}
	fmt.Println(out)
	// Output: host=db.internal port=5432
}
