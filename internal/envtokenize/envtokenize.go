// Package envtokenize splits secret values into ordered tokens
// using a configurable delimiter, enabling structured access to
// composite values such as connection strings or CSV-like fields.
package envtokenize

import (
	"errors"
	"strings"
)

// Tokenizer splits secret values into tokens by a delimiter.
type Tokenizer struct {
	delimiter string
	trim      bool
}

// Result holds the tokens produced from a single secret value.
type Result struct {
	Key    string
	Raw    string
	Tokens []string
}

// New creates a Tokenizer with the given delimiter.
// If trim is true, each token is stripped of surrounding whitespace.
func New(delimiter string, trim bool) (*Tokenizer, error) {
	if delimiter == "" {
		return nil, errors.New("envtokenize: delimiter must not be empty")
	}
	return &Tokenizer{delimiter: delimiter, trim: trim}, nil
}

// Apply tokenizes all values in secrets and returns a map of key → Result.
func (t *Tokenizer) Apply(secrets map[string]string) map[string]Result {
	out := make(map[string]Result, len(secrets))
	for k, v := range secrets {
		tokens := strings.Split(v, t.delimiter)
		if t.trim {
			for i, tok := range tokens {
				tokens[i] = strings.TrimSpace(tok)
			}
		}
		out[k] = Result{Key: k, Raw: v, Tokens: tokens}
	}
	return out
}

// Get returns the n-th token (0-indexed) for the given key.
// Returns an empty string and false if the key is absent or the index is out of range.
func (t *Tokenizer) Get(secrets map[string]string, key string, index int) (string, bool) {
	v, ok := secrets[key]
	if !ok {
		return "", false
	}
	tokens := strings.Split(v, t.delimiter)
	if index < 0 || index >= len(tokens) {
		return "", false
	}
	tok := tokens[index]
	if t.trim {
		tok = strings.TrimSpace(tok)
	}
	return tok, true
}
