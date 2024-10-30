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
		Flags: []cli.Flag{
			&cli.UintFlag{
				Name:  "length",
				Value: s.DefaultTokenLength,
				Usage: "Length of the generated token.",
			},
			&cli.StringFlag{
				Name:  "prefix",
				Value: "",
				Usage: "Prefix for the token.",
			},
		},
		Action: func(ctx *cli.Context) error {
			var (
				namespace *s.Namespace
				account   *s.Account
				err       error
			)

			nsName, accName, token := askForAddTokenDetails(
				ctx.Args().Get(argSetPrefixPositionNamespace),
				ctx.Args().Get(argSetPrefixPositionAccount),
			)

			storage := s.NewFileStorage()
			if err = storage.Prepare(); err != nil {
				return err
			}

			namespace, _ = storage.AddNamespace(&s.Namespace{Name: nsName})

			account, err = namespace.FindAccount(accName)
			if err == nil {
				return CommandError{
					Message: fmt.Sprintf("%s.%s exists", namespace.Name, account.Name),
				}
			}

			account = &s.Account{Name: accName, Token: token, Prefix: ctx.String("prefix"), Length: ctx.Uint("length")}
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
