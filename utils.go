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
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func AskPassword(length int, prompt string) []byte {
	var password []byte = make([]byte, length, length)
	var text string

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

	hash := sha1.New()
	hash.Write([]byte(text))
	h := hash.Sum(nil)
	text = fmt.Sprintf("%x", h)

	copy(password[:], text[0:length])

	return password
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
