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
		Action: func(ctx *cli.Context) (err error) {
			nsName, newName := askForNamespaceRenameDetails(
				ctx.Args().Get(argRenameNamespacePositionNamespace),
				ctx.Args().Get(argRenameNamespacePositionNewName),
			)

			storage := s.NewFileStorage()
			if err = storage.Prepare(); err != nil {
				return err
			}

			defer func() {
				if err != nil {
					return
				}

				err = storage.Save()
			}()

			namespace, err := storage.FindNamespace(nsName)
			if err != nil {
				return resourceNotFoundError(nsName)
			}

			namespace.Name = newName

			return
		},
	}
}

func renameAccountCommand() *cli.Command {
	return &cli.Command{
		Name:      "account",
		ArgsUsage: "[namespace] [account] [new name]",
		Action: func(ctx *cli.Context) (err error) {
			nsName, accName, newName := askForAccountRenameDetails(
				ctx.Args().Get(argRenameAccountPositionNamespace),
				ctx.Args().Get(argRenameAccountPositionAccount),
				ctx.Args().Get(argRenameAccountPositionNewName),
			)

			storage := s.NewFileStorage()
			if err = storage.Prepare(); err != nil {
				return err
			}

			defer func() {
				if err != nil {
					return
				}

				err = storage.Save()
			}()

			namespace, err := storage.FindNamespace(nsName)
			if err != nil {
				return resourceNotFoundError(nsName)
			}

			account, err := namespace.FindAccount(accName)
			if err != nil {
				return resourceNotFoundError(fmt.Sprintf("%s/%s", namespace.Name, accName))
			}

			account.Name = newName

			return
		},
	}
}

func askForNamespaceRenameDetails(namespace, newName string) (string, string) {
	term := terminal.New(os.Stdin, os.Stdout, os.Stderr)

	for len(namespace) < 1 {
		namespace, _ = term.Read("Namespace:")
	}

	for len(newName) < 1 {
		newName, _ = term.Read("New Name:")
	}

	return namespace, newName
}

func askForAccountRenameDetails(namespace, account, newName string) (string, string, string) {
	term := terminal.New(os.Stdin, os.Stdout, os.Stderr)

	for len(namespace) < 1 {
		namespace, _ = term.Read("Namespace:")
	}

	for len(account) < 1 {
		account, _ = term.Read("Account:")
	}

	for len(newName) < 1 {
		newName, _ = term.Read("New Name:")
	}

	return namespace, account, newName
}
