package schema

import (
	"encoding/json"
	"fmt"
	"os"
)

// LoadFile reads a JSON schema file and returns a Schema.
// The JSON format is an array of Rule objects.
func LoadFile(path string) (*Schema, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("schema: open %s: %w", path, err)
	}
	defer f.Close()

	var rules []Rule
	if err := json.NewDecoder(f).Decode(&rules); err != nil {
		return nil, fmt.Errorf("schema: decode %s: %w", path, err)
	}
	return New(rules), nil
}
