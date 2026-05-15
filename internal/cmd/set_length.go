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
		Action: func(ctx *cli.Context) error {
			term := terminal.New(os.Stdin, os.Stdout, os.Stderr)

			nsName, accName, length, err := askForSetLengthDetails(
				term,
				ctx.Args().Get(argSetLengthPositionNamespace),
				ctx.Args().Get(argSetLengthPositionAccount),
				ctx.Args().Get(argSetLengthPositionPrefix),
			)
			if err != nil {
				return err
			}

			storage := prepareStorage(term)

			err = storage.Prepare()
			if err != nil {
				return err
			}

			return executeSetLength(storage, nsName, accName, length)
		},
	}
}

func executeSetLength(storage s.Storage, nsName, accName string, length uint) error {
	namespace, err := storage.FindNamespace(nsName)
	if err != nil {
		return resourceNotFoundError(nsName)
	}

	account, err := namespace.FindAccount(accName)
	if err != nil {
		return resourceNotFoundError(fmt.Sprintf("%s/%s", namespace.Name, accName))
	}

	account.Length = length

	return storage.Save()
}

func askForSetLengthDetails(term terminal.Terminal, namespace, account, length string) (string, string, uint, error) {
	var err error

	for len(namespace) < 1 {
		namespace, err = term.Read("Namespace:")
		if err != nil {
			return "", "", 0, err
		}
	}

	for len(account) < 1 {
		account, err = term.Read("Account:")
		if err != nil {
			return "", "", 0, err
		}
	}

	for {
		for len(length) < 1 {
			length, err = term.Read("Length:")
			if err != nil {
				return "", "", 0, err
			}
		}

		u64Value, parseErr := strconv.ParseUint(length, 10, 32)
		if parseErr == nil {
			return namespace, account, uint(u64Value), nil
		}

		length = ""
	}
}
