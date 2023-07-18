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
			&cli.BoolFlag{
				Name:  "clear",
				Value: false,
				Usage: "Clear prefix from account.",
			},
		},
		Action: func(ctx *cli.Context) (err error) {
			var (
				namespace *s.Namespace
				account   *s.Account
			)

			nsName, accName, prefix := askForSetPrefixDetails(
				ctx.Args().Get(argPositionNamespace),
				ctx.Args().Get(argPositionAccount),
				ctx.Args().Get(argPositionPrefix),
				ctx.Bool("clear"),
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

			account.Prefix = prefix

			return nil
		},
	}
}

func askForSetPrefixDetails(namespace, account, prefix string, clear bool) (string, string, string) {
	term := terminal.New(os.Stdin, os.Stdout, os.Stderr)

	for len(namespace) < 1 {
		namespace, _ = term.Read("Namespace:")
	}

	for len(account) < 1 {
		account, _ = term.Read("Account:")
	}

	if clear {
		return namespace, account, ""
	}

	for len(prefix) < 1 {
		prefix, _ = term.Read("Prefix:")
	}

	return namespace, account, prefix
}
