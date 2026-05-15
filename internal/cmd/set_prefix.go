package cmd

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	s "github.com/yitsushi/totp-cli/internal/storage"
	"github.com/yitsushi/totp-cli/internal/terminal"
)

// SetPrefixCommand is the set-prefix subcommand.
func SetPrefixCommand() *cli.Command {
	return &cli.Command{
		Name:      "set-prefix",
		Usage:     "Set prefix for a token.",
		ArgsUsage: "[namespace] [account] [prefix]",
		Flags: []cli.Flag{
			flagClearPrefix(),
		},
		Action: func(ctx *cli.Context) error {
			term := terminal.New(os.Stdin, os.Stdout, os.Stderr)

			nsName, accName, prefix, err := askForSetPrefixDetails(
				term,
				ctx.Args().Get(argSetPrefixPositionNamespace),
				ctx.Args().Get(argSetPrefixPositionAccount),
				ctx.Args().Get(argSetPrefixPositionPrefix),
				ctx.Bool("clear"),
			)
			if err != nil {
				return err
			}

			storage := prepareStorage(term)

			err = storage.Prepare()
			if err != nil {
				return err
			}

			return executeSetPrefix(storage, nsName, accName, prefix)
		},
	}
}

func executeSetPrefix(storage s.Storage, nsName, accName, prefix string) error {
	namespace, err := storage.FindNamespace(nsName)
	if err != nil {
		return resourceNotFoundError(nsName)
	}

	account, err := namespace.FindAccount(accName)
	if err != nil {
		return resourceNotFoundError(fmt.Sprintf("%s/%s", namespace.Name, accName))
	}

	account.Prefix = prefix

	return storage.Save()
}

func askForSetPrefixDetails(
	term terminal.Terminal, namespace, account, prefix string, isClear bool,
) (string, string, string, error) {
	var err error

	for len(namespace) < 1 {
		namespace, err = term.Read("Namespace:")
		if err != nil {
			return "", "", "", err
		}
	}

	for len(account) < 1 {
		account, err = term.Read("Account:")
		if err != nil {
			return "", "", "", err
		}
	}

	if isClear {
		return namespace, account, "", nil
	}

	for len(prefix) < 1 {
		prefix, err = term.Read("Prefix:")
		if err != nil {
			return "", "", "", err
		}
	}

	return namespace, account, prefix, nil
}
