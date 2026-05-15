package cmd

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	s "github.com/yitsushi/totp-cli/internal/storage"
	"github.com/yitsushi/totp-cli/internal/terminal"
)

// DeleteCommand is subcommand to delete an account or a namespace.
func DeleteCommand() *cli.Command {
	return &cli.Command{
		Name:      "delete",
		Usage:     "Delete an account or a whole namespace.",
		ArgsUsage: "<namespace> [account]",
		Action: func(ctx *cli.Context) error {
			namespaceName := ctx.Args().Get(argSetPrefixPositionNamespace)
			if len(namespaceName) < 1 {
				return CommandError{Message: errMsgNamespaceNotDefined}
			}

			term := terminal.New(os.Stdin, os.Stdout, os.Stderr)
			storage := prepareStorage(term)

			err := storage.Prepare()
			if err != nil {
				return err
			}

			return executeDelete(storage, term, namespaceName, ctx.Args().Get(argSetPrefixPositionAccount))
		},
	}
}

func executeDelete(storage s.Storage, term terminal.Terminal, nsName, accName string) error {
	namespace, err := storage.FindNamespace(nsName)
	if err != nil {
		return err
	}

	if accName != "" {
		account, err := namespace.FindAccount(accName)
		if err != nil {
			return err
		}

		fmt.Printf("You want to delete '%s.%s' account.\n", namespace.Name, account.Name)

		if term.Confirm("Are you sure?") {
			namespace.DeleteAccount(account)
		}

		return storage.Save()
	}

	fmt.Printf("You want to delete '%s' namespace with %d accounts.\n", namespace.Name, len(namespace.Accounts))

	for _, account := range namespace.Accounts {
		fmt.Printf(" - %s.%s\n", namespace.Name, account.Name)
	}

	if term.Confirm("Are you sure?") {
		storage.DeleteNamespace(namespace)
	}

	return storage.Save()
}
