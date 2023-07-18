package cmd

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	s "github.com/yitsushi/totp-cli/internal/storage"
	"github.com/yitsushi/totp-cli/internal/terminal"
)

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

			if !CheckPasswordConfirm([]byte(newPasswordIn), []byte(newPasswordConfirmIn)) {
				return fmt.Errorf("New Password and the confirm mismatch!")
			}

			storage.Password = newPasswordIn

			if err = storage.Save(); err != nil {
				return err
			}

			return nil
		},
	}
}
