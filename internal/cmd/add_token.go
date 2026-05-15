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
			flagLength(),
			flagPrefix(),
			flagAlgorithm(),
			flagTimePeriod(),
		},
		Action: func(ctx *cli.Context) error {
			term := terminal.New(os.Stdin, os.Stdout, os.Stderr)

			nsName, accName, token, err := askForAddTokenDetails(
				term,
				ctx.Args().Get(argSetPrefixPositionNamespace),
				ctx.Args().Get(argSetPrefixPositionAccount),
			)
			if err != nil {
				return err
			}

			storage := prepareStorage(term)

			err = storage.Prepare()
			if err != nil {
				return err
			}

			return executeAddToken(storage, nsName, accName, token, AccountOptions{
				Prefix:     ctx.String("prefix"),
				Length:     ctx.Uint("length"),
				Algorithm:  ctx.String("algorithm"),
				TimePeriod: ctx.Int64("time-period"),
			})
		},
	}
}

func executeAddToken(storage s.Storage, nsName, accName, token string, opts AccountOptions) error {
	namespace, _ := storage.AddNamespace(&s.Namespace{Name: nsName})

	_, err := namespace.FindAccount(accName)
	if err == nil {
		return CommandError{
			Message: fmt.Sprintf("%s.%s exists", namespace.Name, accName),
		}
	}

	namespace.Accounts = append(namespace.Accounts, &s.Account{
		Name:       accName,
		Token:      token,
		Prefix:     opts.Prefix,
		Length:     opts.Length,
		Algorithm:  opts.Algorithm,
		TimePeriod: opts.TimePeriod,
	})

	return storage.Save()
}

func askForAddTokenDetails(term terminal.Terminal, namespace, account string) (string, string, string, error) {
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

	token, err := term.Read("Token:")

	return namespace, account, token, err
}
