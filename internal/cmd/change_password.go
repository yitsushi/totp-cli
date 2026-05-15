package cmd

import (
	"os"

	"github.com/urfave/cli/v2"
	"github.com/yitsushi/totp-cli/internal/security"
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
			term := terminal.New(os.Stdin, os.Stdout, os.Stderr)
			storage := prepareStorage(term)

			err := storage.Prepare()
			if err != nil {
				return err
			}

			newPassword, err := askForNewPassword(term)
			if err != nil {
				return err
			}

			return executeChangePassword(storage, newPassword)
		},
	}
}

func executeChangePassword(storage s.Storage, newPassword string) error {
	storage.SetPassword(newPassword)

	return storage.Save()
}

func askForNewPassword(term terminal.Terminal) (string, error) {
	newPassword, err := term.Hidden("New Password:")
	if err != nil {
		return "", err
	}

	confirm, err := term.Hidden("Again:")
	if err != nil {
		return "", err
	}

	if !security.CheckPasswordConfirm([]byte(newPassword), []byte(confirm)) {
		return "", CommandError{Message: "new password and the confirm mismatch"}
	}

	return newPassword, nil
}
