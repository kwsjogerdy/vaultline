package env

import (
	"fmt"
	"os"
	"strings"
)

// Writer handles writing secrets to .env files.
type Writer struct {
	FilePath string
}

// NewWriter creates a new Writer for the given file path.
func NewWriter(filePath string) *Writer {
	return &Writer{FilePath: filePath}
}

// Write writes the provided secrets map to the .env file.
// Existing file content is overwritten.
func (w *Writer) Write(secrets map[string]string) error {
	f, err := os.OpenFile(w.FilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("env writer: failed to open file %q: %w", w.FilePath, err)
	}
	defer f.Close()

	for key, value := range secrets {
		line := fmt.Sprintf("%s=%s\n", sanitizeKey(key), escapeValue(value))
		if _, err := f.WriteString(line); err != nil {
			return fmt.Errorf("env writer: failed to write key %q: %w", key, err)
		}
	}
	return nil
}

// sanitizeKey uppercases and trims whitespace from a key.
func sanitizeKey(key string) string {
	return strings.ToUpper(strings.TrimSpace(key))
}

// escapeValue wraps values containing spaces or special chars in double quotes.
func escapeValue(value string) string {
	if strings.ContainsAny(value, " \t\n#") {
		escaped := strings.ReplaceAll(value, `"`, `\"`)
		return fmt.Sprintf(`"%s"`, escaped)
	}
	return value
}
