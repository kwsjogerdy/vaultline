package prompt

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// Prompter handles interactive user confirmation.
type Prompter struct {
	in  io.Reader
	out io.Writer
}

// New returns a Prompter using stdin/stdout.
func New() *Prompter {
	return &Prompter{in: os.Stdin, out: os.Stdout}
}

// NewWithIO returns a Prompter with custom IO (useful for testing).
func NewWithIO(in io.Reader, out io.Writer) *Prompter {
	return &Prompter{in: in, out: out}
}

// Confirm asks the user a yes/no question and returns true if they answer yes.
func (p *Prompter) Confirm(question string) (bool, error) {
	fmt.Fprintf(p.out, "%s [y/N]: ", question)
	scanner := bufio.NewScanner(p.in)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return false, fmt.Errorf("reading input: %w", err)
		}
		return false, nil
	}
	answer := strings.TrimSpace(strings.ToLower(scanner.Text()))
	return answer == "y" || answer == "yes", nil
}

// ConfirmDiff displays a summary of changes and asks the user to confirm.
func (p *Prompter) ConfirmDiff(added, removed, changed int) (bool, error) {
	fmt.Fprintf(p.out, "\nPending changes:\n")
	fmt.Fprintf(p.out, "  + %d added\n", added)
	fmt.Fprintf(p.out, "  - %d removed\n", removed)
	fmt.Fprintf(p.out, "  ~ %d changed\n", changed)
	fmt.Fprintf(p.out, "\n")
	return p.Confirm("Apply these changes?")
}
