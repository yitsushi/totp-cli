package cmd

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	s "github.com/yitsushi/totp-cli/internal/storage"
	"github.com/yitsushi/totp-cli/internal/terminal"
)

// AddTokenCommand is the add-token subcommand.
func AddTokenCommand() *cli.Command {
	return &cli.Command{
		Name:      "add-token",
		Aliases:   []string{"add"},
		Usage:     "Add new token.",
		ArgsUsage: "[namespace] [account]",
		Action: func(ctx *cli.Context) error {
			var (
				namespace *s.Namespace
				account   *s.Account
				err       error
			)

			nsName, accName, token := askForAddTokenDetails(
				ctx.Args().Get(argPositionNamespace),
				ctx.Args().Get(argPositionAccount),
			)

			storage, err := s.PrepareStorage()
			if err != nil {
				return err
			}

			namespace, err = storage.FindNamespace(nsName)
			if err != nil {
				namespace = &s.Namespace{Name: nsName}
				storage.Namespaces = append(storage.Namespaces, namespace)
			}

			account, err = namespace.FindAccount(accName)
			if err == nil {
				return CommandError{
					Message: fmt.Sprintf("%s.%s exists", namespace.Name, account.Name),
				}
			}

			account = &s.Account{Name: accName, Token: token}
			namespace.Accounts = append(namespace.Accounts, account)

			err = storage.Save()
			if err != nil {
				return err
			}

			return nil
		},
	}
}

func askForAddTokenDetails(namespace, account string) (string, string, string) {
	term := terminal.New(os.Stdin, os.Stdout, os.Stderr)

	for len(namespace) < 1 {
		namespace, _ = term.Read("Namespace:")
	}

	for len(account) < 1 {
		account, _ = term.Read("Account:")
	}

	token, _ := term.Read("Token:")

	return namespace, account, token
}
