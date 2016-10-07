package command

import (
	"fmt"

	s "github.com/Yitsushi/totp-cli/storage"
	"github.com/Yitsushi/totp-cli/util"
)

type ChangePassword struct {
}

func (c *ChangePassword) Execute() {
	storage := s.PrepareStorage()
	newPassword := util.AskPassword(32, "New Password")
	newPasswordConfirm := util.AskPassword(32, "Again")

	if !util.CheckPasswordConfirm(newPassword, newPasswordConfirm) {
		fmt.Println("New Password and the confirm mismatch!")
		return
	}

	storage.Password = newPassword
	storage.Save()
}

func (c *ChangePassword) ArgumentDescription() string {
	return ""
}

func (c *ChangePassword) Description() string {
	return "Change password"
}

func (c *ChangePassword) Help() string {
	return ""
}

func (c *ChangePassword) Examples() []string {
	return []string{""}
}
