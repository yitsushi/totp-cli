package cmd

import (
	"fmt"
	"os"

	"github.com/yitsushi/go-commander"

	s "github.com/yitsushi/totp-cli/internal/storage"
	"github.com/yitsushi/totp-cli/internal/util"
)

const (
	askPasswordLength = 32
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

	newPassword := util.AskPassword(askPasswordLength, "New Password")
	newPasswordConfirm := util.AskPassword(askPasswordLength, "Again")

	if !util.CheckPasswordConfirm(newPassword, newPasswordConfirm) {
		fmt.Println("New Password and the confirm mismatch!")
		os.Exit(1)
	}

	storage.Password = newPassword

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
