package main

import (
	"bufio"
	"crypto/sha1"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func AskPIN(length int) []byte {
	var pin []byte = make([]byte, length, length)
	var text string

	text = os.Getenv("PIN")

	if len(text) < 1 {
		os.Stderr.Write([]byte("PIN: "))
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

	hash := sha1.New()
	hash.Write([]byte(text))
	h := hash.Sum(nil)
	text = fmt.Sprintf("%x", h)

	copy(pin[:], text[0:length])

	return pin
}

func Ask(prompt string) (text string) {
	fmt.Printf("%s: ", prompt)
	reader := bufio.NewReader(os.Stdin)
	text, _ = reader.ReadString('\n')
	text = strings.TrimSpace(text)

	return
}

func Confirm(prompt string) bool {
	fmt.Printf("%s ", prompt)
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = strings.ToLower(strings.TrimSpace(text))

	return (text == "y" || text == "yes" || text == "sure")
}
