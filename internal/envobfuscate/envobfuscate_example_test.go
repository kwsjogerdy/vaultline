package envobfuscate_test

import (
	"fmt"

	"github.com/vaultline/vaultline/internal/envobfuscate"
)

func ExampleObfuscator_Obfuscate() {
	o, err := envobfuscate.New("my-shared-key")
	if err != nil {
		panic(err)
	}
	secrets := map[string]string{
		"TOKEN": "s3cr3t",
	}
	obf := o.Obfuscate(secrets)
	fmt.Println(envobfuscate.IsObfuscated(obf["TOKEN"]))

	restored, _ := o.Deobfuscate(obf)
	fmt.Println(restored["TOKEN"])
	// Output:
	// true
	// s3cr3t
}
