package util

import (
	"bufio"

	//nolint:gosec // yolo?
	"crypto/sha1"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

// AskPassword asks password from the user and hides
// the input.
func AskPassword(length int, prompt string) []byte {
	var text string

	password := make([]byte, length)

	if prompt == "" {
		prompt = "Password"
	}

	text = os.Getenv("Password")

	if len(text) < 1 {
		prompt = fmt.Sprintf("%s: ", prompt)
		os.Stderr.Write([]byte(prompt))
		stdin := int(os.Stdin.Fd())

		if terminal.IsTerminal(stdin) {
			var lineBytes []byte

			lineBytes, _ = terminal.ReadPassword(stdin)
			text = string(lineBytes)
		} else {
			reader := bufio.NewReader(os.Stdin)
			text, _ = reader.ReadString('\n')
		}

		os.Stderr.Write([]byte("***\n"))

		text = strings.TrimSpace(text)
	}

	hash := sha1.New() //nolint:gosec // yolo?
	_, _ = hash.Write([]byte(text))
	h := hash.Sum(nil)
	text = fmt.Sprintf("%x", h)

	copy(password[:], text[0:length]) //nolint:gocritic // intentional

	return password
}

// Ask can be used to get some input from the user.
// The user input will not be hidden (not secure).
func Ask(prompt string) (text string) {
	fmt.Printf("%s: ", prompt)

	reader := bufio.NewReader(os.Stdin)
	text, _ = reader.ReadString('\n')
	text = strings.TrimSpace(text)

	return
}

// Read stdin without prompt..
func Read() (text string) {
	reader := bufio.NewReader(os.Stdin)
	text, _ = reader.ReadString('\n')
	text = strings.TrimSpace(text)

	return
}

// Confirm ask something from the user
// Acceptable true answers: yes, y, sure
// Everything else will be false.
func Confirm(prompt string) bool {
	fmt.Printf("%s ", prompt)

	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = strings.ToLower(strings.TrimSpace(text))

	return (text == "y" || text == "yes" || text == "sure")
}

// CheckPasswordConfirm checks two byte array if the content is the same.
func CheckPasswordConfirm(password, confirm []byte) bool {
	if password == nil && confirm == nil {
		return true
	}

	if password == nil || confirm == nil {
		return false
	}

	if len(password) != len(confirm) {
		return false
	}

	for i := range password {
		if password[i] != confirm[i] {
			return false
		}
	}

	return true
}
