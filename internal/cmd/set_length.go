package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/urfave/cli/v2"

	s "github.com/yitsushi/totp-cli/internal/storage"
	"github.com/yitsushi/totp-cli/internal/terminal"
)

// SetLengthCommand is the set-prefix subcommand.
func SetLengthCommand() *cli.Command {
	return &cli.Command{
		Name:      "set-length",
		Usage:     "Set length for a token.",
		ArgsUsage: "[namespace] [account] [length]",
		Action: func(ctx *cli.Context) (err error) {
			var (
				namespace *s.Namespace
				account   *s.Account
			)

			nsName, accName, length := askForSetLengthDetails(
				ctx.Args().Get(argSetLengthPositionNamespace),
				ctx.Args().Get(argSetLengthPositionAccount),
				ctx.Args().Get(argSetLengthPositionPrefix),
			)

			storage, err := s.PrepareStorage()
			if err != nil {
				return
			}

			defer func() {
				if err != nil {
					return
				}

				err = storage.Save()
			}()

			namespace, err = storage.FindNamespace(nsName)
			if err != nil {
				return CommandError{Message: fmt.Sprintf("%s does not exist", nsName)}
			}

			account, err = namespace.FindAccount(accName)
			if err != nil {
				return CommandError{Message: fmt.Sprintf("%s/%s does not exist", namespace.Name, accName)}
			}

			account.Length = length

			return nil
		},
	}
}

func askForSetLengthDetails(namespace, account, length string) (string, string, uint) {
	term := terminal.New(os.Stdin, os.Stdout, os.Stderr)

	for len(namespace) < 1 {
		namespace, _ = term.Read("Namespace:")
	}

	for len(account) < 1 {
		account, _ = term.Read("Account:")
	}

	for len(length) < 1 {
		length, _ = term.Read("Length:")
	}

	u64Value, _ := strconv.ParseUint(length, 10, 32)

	return namespace, account, uint(u64Value)
}
