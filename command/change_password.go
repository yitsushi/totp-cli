package command

import (
	"github.com/yitsushi/go-commander"
	s "github.com/yitsushi/totp-cli/storage"
	"github.com/yitsushi/totp-cli/util"
)

const (
	askPasswordLength = 32
)

// ChangePassword structure is the representation of the change-password command.
type ChangePassword struct {
}

// Execute is the main function. It will be called on change-password command.
func (c *ChangePassword) Execute(opts *commander.CommandHelper) {
	storage := s.PrepareStorage()
	newPassword := util.AskPassword(askPasswordLength, "New Password")
	newPasswordConfirm := util.AskPassword(askPasswordLength, "Again")

	if !util.CheckPasswordConfirm(newPassword, newPasswordConfirm) {
		panic("New Password and the confirm mismatch!")
	}

	storage.Password = newPassword
	storage.Save()
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
