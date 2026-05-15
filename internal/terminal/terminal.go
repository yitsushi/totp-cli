package terminal

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"
)

// Terminal is the interface for interacting with the user.
type Terminal interface {
	Read(prompt string) (string, error)
	Confirm(prompt string) bool
	Hidden(prompt string) (string, error)
}

// ConcreteTerminal implements Terminal using an input reader and output writers.
type ConcreteTerminal struct {
	Input       io.Reader
	Output      io.Writer
	ErrorOutput io.Writer
	reader      *bufio.Reader
}

// New terminal instance.
func New(in io.Reader, out io.Writer, errorOut io.Writer) *ConcreteTerminal {
	return &ConcreteTerminal{
		Input:       in,
		Output:      out,
		ErrorOutput: errorOut,
	}
}

// Read text from input with optional prompt.
func (t *ConcreteTerminal) Read(prompt string) (string, error) {
	if prompt != "" {
		_, _ = fmt.Fprintf(t.Output, "%s ", prompt)
	}

	if t.reader == nil {
		t.reader = bufio.NewReader(t.Input)
	}

	text, readErr := t.reader.ReadString('\n')
	if readErr != nil {
		return text, fmt.Errorf("error reading from input: %w", readErr)
	}

	text = strings.TrimSpace(text)

	return text, nil
}

// Confirm asks the user for confirmation.
// If the answer is y, yes, or sure, the used confirmed,
// otherwise not.
func (t *ConcreteTerminal) Confirm(prompt string) bool {
	answer, err := t.Read(prompt)
	if err != nil {
		return false
	}

	answer = strings.ToLower(strings.TrimSpace(answer))

	return answer == "y" || answer == "yes" || answer == "sure"
}

// Hidden reads from Input, but typed characters are hidden.
// Good for passwords, tokens, or other sensitive information.
func (t *ConcreteTerminal) Hidden(prompt string) (string, error) {
	var (
		text string
		err  error
	)

	if prompt != "" {
		_, _ = fmt.Fprintf(t.ErrorOutput, "%s ", prompt)
	}

	in, inIsFile := t.Input.(*os.File)

	if inIsFile && term.IsTerminal(int(in.Fd())) {
		var lineBytes []byte

		lineBytes, err = term.ReadPassword(int(in.Fd()))
		text = string(lineBytes)
	} else {
		text, err = t.Read("")
	}

	_, _ = fmt.Fprintln(t.ErrorOutput, "***")

	return strings.TrimSpace(text), err
}
