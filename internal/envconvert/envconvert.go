// Package envconvert converts secrets between different value formats.
package envconvert

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"
)

// Format represents a conversion target format.
type Format string

const (
	FormatBase64    Format = "base64"
	FormatURLEncode Format = "urlencode"
	FormatUppercase Format = "upper"
	FormatLowercase Format = "lower"
	FormatTrimSpace Format = "trim"
)

// Result holds the original and converted value for a key.
type Result struct {
	Key      string
	Original string
	Converted string
}

// Converter applies a format transformation to secrets.
type Converter struct {
	format Format
}

// New creates a Converter for the given format.
func New(format Format) (*Converter, error) {
	switch format {
	case FormatBase64, FormatURLEncode, FormatUppercase, FormatLowercase, FormatTrimSpace:
		return &Converter{format: format}, nil
	default:
		return nil, fmt.Errorf("unknown format %q", format)
	}
}

// Apply converts all values in secrets and returns a Result slice and updated map.
func (c *Converter) Apply(secrets map[string]string) ([]Result, map[string]string) {
	out := make(map[string]string, len(secrets))
	results := make([]Result, 0, len(secrets))
	for k, v := range secrets {
		converted := c.convert(v)
		out[k] = converted
		results = append(results, Result{Key: k, Original: v, Converted: converted})
	}
	return results, out
}

// Format returns the format string associated with this Converter.
func (c *Converter) Format() Format {
	return c.format
}

// SupportedFormats returns a slice of all valid format values.
func SupportedFormats() []Format {
	return []Format{
		FormatBase64,
		FormatURLEncode,
		FormatUppercase,
		FormatLowercase,
		FormatTrimSpace,
	}
}

func (c *Converter) convert(v string) string {
	switch c.format {
	case FormatBase64:
		return base64.StdEncoding.EncodeToString([]byte(v))
	case FormatURLEncode:
		return url.QueryEscape(v)
	case FormatUppercase:
		return strings.ToUpper(v)
	case FormatLowercase:
		return strings.ToLower(v)
	case FormatTrimSpace:
		return strings.TrimSpace(v)
	default:
		return v
	}
}
