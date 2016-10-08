package command

import (
	s "github.com/Yitsushi/totp-cli/storage"
	"github.com/Yitsushi/totp-cli/util"
)

// ChangePassword stucture is the representation of the change-password command
type ChangePassword struct {
}

// Execute is the main function. It will be called on change-password command
func (c *ChangePassword) Execute() {
	storage := s.PrepareStorage()
	newPassword := util.AskPassword(32, "New Password")
	newPasswordConfirm := util.AskPassword(32, "Again")

	if !util.CheckPasswordConfirm(newPassword, newPasswordConfirm) {
		panic("New Password and the confirm mismatch!")
	}

	storage.Password = newPassword
	storage.Save()
}

// ArgumentDescription descripts the required and potential arguments
func (c *ChangePassword) ArgumentDescription() string {
	return ""
}

// Description will be displayed as Description (woooo) in the general help
func (c *ChangePassword) Description() string {
	return "Change password"
}

// Help is a general (human readable) command specific (long) help
func (c *ChangePassword) Help() string {
	return ""
}

// Examples lists a few example as array. Will be used in the command specific help
func (c *ChangePassword) Examples() []string {
	return []string{""}
}
