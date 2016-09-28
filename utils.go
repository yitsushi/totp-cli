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

func AskPIN(length int, prompt string) []byte {
	var pin []byte = make([]byte, length, length)
	var text string

	if prompt == "" {
		prompt = "PIN"
	}

	text = os.Getenv("PIN")

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

func CheckPINConfirm(pin, confirm []byte) bool {
	if pin == nil && confirm == nil {
		return true
	}

	if pin == nil || confirm == nil {
		return false
	}

	if len(pin) != len(confirm) {
		return false
	}

	for i := range pin {
		if pin[i] != confirm[i] {
			return false
		}
	}

	return true
}
