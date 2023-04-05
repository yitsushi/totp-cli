package cmd

import (
	"fmt"
	"os"

	"github.com/yitsushi/go-commander"

	s "github.com/yitsushi/totp-cli/internal/storage"
	"github.com/yitsushi/totp-cli/internal/terminal"
)

// ChangePassword structure is the representation of the change-password command.
type ChangePassword struct{}

// Execute is the main function. It will be called on change-password command.
func (c *ChangePassword) Execute(opts *commander.CommandHelper) {
	storage, err := s.PrepareStorage()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	term := terminal.New(os.Stdin, os.Stdout, os.Stderr)

	newPasswordIn, err := term.Hidden("New Password:")
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	newPasswordConfirmIn, err := term.Hidden("Again:")
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	if !CheckPasswordConfirm([]byte(newPasswordIn), []byte(newPasswordConfirmIn)) {
		fmt.Println("New Password and the confirm mismatch!")
		os.Exit(1)
	}

	storage.Password = newPasswordIn

	err = storage.Save()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
}

// NewChangePassword create a new ChangePassword command.
func NewChangePassword(appName string) *commander.CommandWrapper {
	return &commander.CommandWrapper{
		Handler: &ChangePassword{},
		Help: &commander.CommandDescriptor{
			Name:             "change-password",
			ShortDescription: "Change password",
		},
	}
}
