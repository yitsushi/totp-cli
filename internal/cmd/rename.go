package cmd

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	s "github.com/yitsushi/totp-cli/internal/storage"
	"github.com/yitsushi/totp-cli/internal/terminal"
)

const (
	argRenameNamespacePositionNamespace = 0
	argRenameNamespacePositionNewName   = 1
	argRenameAccountPositionNamespace   = 0
	argRenameAccountPositionAccount     = 1
	argRenameAccountPositionNewName     = 2
)

// RenameCommand is a command to rename existing namespace or account.
func RenameCommand() *cli.Command {
	return &cli.Command{
		Name:  "rename",
		Usage: "Rename an account or namespace",
		Subcommands: []*cli.Command{
			renameAccountCommand(),
			renameNamespaceCommand(),
		},
	}
}

func renameNamespaceCommand() *cli.Command {
	return &cli.Command{
		Name:      "namespace",
		ArgsUsage: "[namespace] [new name]",
		Action: func(ctx *cli.Context) error {
			term := terminal.New(os.Stdin, os.Stdout, os.Stderr)

			nsName, newName, err := askForNamespaceRenameDetails(
				term,
				ctx.Args().Get(argRenameNamespacePositionNamespace),
				ctx.Args().Get(argRenameNamespacePositionNewName),
			)
			if err != nil {
				return err
			}

			storage := prepareStorage(term)

			err = storage.Prepare()
			if err != nil {
				return err
			}

			return executeRenameNamespace(storage, nsName, newName)
		},
	}
}

func renameAccountCommand() *cli.Command {
	return &cli.Command{
		Name:      "account",
		ArgsUsage: "[namespace] [account] [new name]",
		Action: func(ctx *cli.Context) error {
			term := terminal.New(os.Stdin, os.Stdout, os.Stderr)

			nsName, accName, newName, err := askForAccountRenameDetails(
				term,
				ctx.Args().Get(argRenameAccountPositionNamespace),
				ctx.Args().Get(argRenameAccountPositionAccount),
				ctx.Args().Get(argRenameAccountPositionNewName),
			)
			if err != nil {
				return err
			}

			storage := prepareStorage(term)

			err = storage.Prepare()
			if err != nil {
				return err
			}

			return executeRenameAccount(storage, nsName, accName, newName)
		},
	}
}

func executeRenameNamespace(storage s.Storage, nsName, newName string) error {
	namespace, err := storage.FindNamespace(nsName)
	if err != nil {
		return resourceNotFoundError(nsName)
	}

	namespace.Name = newName

	return storage.Save()
}

func executeRenameAccount(storage s.Storage, nsName, accName, newName string) error {
	namespace, err := storage.FindNamespace(nsName)
	if err != nil {
		return resourceNotFoundError(nsName)
	}

	account, err := namespace.FindAccount(accName)
	if err != nil {
		return resourceNotFoundError(fmt.Sprintf("%s/%s", namespace.Name, accName))
	}

	account.Name = newName

	return storage.Save()
}

func askForNamespaceRenameDetails(term terminal.Terminal, namespace, newName string) (string, string, error) {
	var err error

	for len(namespace) < 1 {
		namespace, err = term.Read("Namespace:")
		if err != nil {
			return "", "", err
		}
	}

	for len(newName) < 1 {
		newName, err = term.Read("New Name:")
		if err != nil {
			return "", "", err
		}
	}

	return namespace, newName, nil
}

func askForAccountRenameDetails(
	term terminal.Terminal, namespace, account, newName string,
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

	for len(newName) < 1 {
		newName, err = term.Read("New Name:")
		if err != nil {
			return "", "", "", err
		}
	}

	return namespace, account, newName, nil
}
