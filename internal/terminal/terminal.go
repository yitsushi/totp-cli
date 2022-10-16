package terminal

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"
)

// Terminal represents the terminal with input, output and error output.
type Terminal struct {
	Input       io.Reader
	Output      io.Writer
	ErrorOutput io.Writer
}

// New terminal instance.
func New(in io.Reader, out io.Writer, errorOut io.Writer) Terminal {
	return Terminal{
		Input:       in,
		Output:      out,
		ErrorOutput: errorOut,
	}
}

// Read text from input with optinal prompt.
func (t Terminal) Read(prompt string) (string, error) {
	if prompt != "" {
		fmt.Fprintf(t.Output, "%s ", prompt)
	}

	reader := bufio.NewReader(t.Input)

	text, readErr := reader.ReadString('\n')
	if readErr != nil {
		return "", fmt.Errorf("error reading from input: %w", readErr)
	}

	text = strings.TrimSpace(text)

	return text, nil
}

// Confirm asks the user for confirmation.
// If the answer is y, yes, or sure, the used confirmed,
// otherwise not.
func (t Terminal) Confirm(prompt string) bool {
	answer, err := t.Read(prompt)
	if err != nil {
		return false
	}

	answer = strings.ToLower(strings.TrimSpace(answer))

	return answer == "y" || answer == "yes" || answer == "sure"
}

// Hidden reads from Input, but typed characters are hidden.
// Good for passwords, tokens, or other sensitive information.
func (t Terminal) Hidden(prompt string) (string, error) {
	var (
		text string
		err  error
	)

	if prompt != "" {
		fmt.Fprintf(t.ErrorOutput, "%s ", prompt)
	}

	in, inIsFile := t.Input.(*os.File)

	if inIsFile && term.IsTerminal(int(in.Fd())) {
		var lineBytes []byte

		lineBytes, err = term.ReadPassword(int(in.Fd()))
		text = string(lineBytes)
	} else {
		text, err = t.Read("")
	}

	fmt.Fprintln(t.ErrorOutput, "***")

	return strings.TrimSpace(text), err
}
