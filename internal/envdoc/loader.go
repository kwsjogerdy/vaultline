package envdoc

import (
	"encoding/json"
	"fmt"
	"os"
)

// LoadFile reads a JSON file mapping key -> annotation string and returns a Doc.
// The JSON format is: {"KEY": "description|example|required"}
func LoadFile(path string) (*Doc, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("envdoc: open %s: %w", path, err)
	}
	defer f.Close()

	var annotations map[string]string
	if err := json.NewDecoder(f).Decode(&annotations); err != nil {
		return nil, fmt.Errorf("envdoc: decode %s: %w", path, err)
	}
	return New(annotations)
}
