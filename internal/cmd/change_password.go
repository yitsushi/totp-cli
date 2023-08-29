package cmd

import (
	"github.com/yitsushi/totp-cli/internal/security"
	"os"

	"github.com/urfave/cli/v2"

	s "github.com/yitsushi/totp-cli/internal/storage"
	"github.com/yitsushi/totp-cli/internal/terminal"
)

// ChangePasswordCommand is the change-password subcommand.
func ChangePasswordCommand() *cli.Command {
	return &cli.Command{
		Name:      "change-password",
		Usage:     "Change password.",
		ArgsUsage: "",
		Action: func(_ *cli.Context) error {
			var (
				err                  error
				storage              *s.Storage
				newPasswordIn        string
				newPasswordConfirmIn string
			)

			if storage, err = s.PrepareStorage(); err != nil {
				return err
			}

			term := terminal.New(os.Stdin, os.Stdout, os.Stderr)

			if newPasswordIn, err = term.Hidden("New Password:"); err != nil {
				return err
			}

			if newPasswordConfirmIn, err = term.Hidden("Again:"); err != nil {
				return err
			}

			if !security.CheckPasswordConfirm([]byte(newPasswordIn), []byte(newPasswordConfirmIn)) {
				return CommandError{Message: "new password and the confirm mismatch"}
			}

			storage.Password = newPasswordIn

			return storage.Save()
		},
	}
}
