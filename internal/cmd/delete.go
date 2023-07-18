package cmd

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	s "github.com/yitsushi/totp-cli/internal/storage"
	"github.com/yitsushi/totp-cli/internal/terminal"
)

// DeleteCommand is the delete subcommand.
func DeleteCommand() *cli.Command {
	return &cli.Command{
		Name:      "delete",
		Usage:     "Delete an account or a whole namespace.",
		ArgsUsage: "<namespace> [account]",
		Action: func(ctx *cli.Context) (err error) {
			namespaceName := ctx.Args().Get(argPositionNamespace)
			if len(namespaceName) < 1 {
				return CommandError{Message: "namespace is not defined"}
			}

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

			var (
				namespace *s.Namespace
				account   *s.Account
			)

			if namespace, err = storage.FindNamespace(namespaceName); err != nil {
				return
			}

			term := terminal.New(os.Stdin, os.Stdout, os.Stderr)

			if accountName := ctx.Args().Get(argPositionAccount); accountName != "" {
				if account, err = namespace.FindAccount(accountName); err != nil {
					return
				}

				fmt.Printf("You want to delete '%s.%s' account.\n", namespace.Name, account.Name)

				if term.Confirm("Are you sure?") {
					namespace.DeleteAccount(account)
				}

				return
			}

			fmt.Printf("You want to delete '%s' namespace with %d accounts.\n", namespace.Name, len(namespace.Accounts))
			for _, account := range namespace.Accounts {
				fmt.Printf(" - %s.%s\n", namespace.Name, account.Name)
			}

			if term.Confirm("Are you sure?") {
				storage.DeleteNamespace(namespace)
			}

			return
		},
	}
}
