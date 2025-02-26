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
		Action: func(ctx *cli.Context) (err error) {
			var (
				namespace *s.Namespace
				account   *s.Account
			)

			nsName, accName, prefix := askForSetPrefixDetails(
				ctx.Args().Get(argSetPrefixPositionNamespace),
				ctx.Args().Get(argSetPrefixPositionAccount),
				ctx.Args().Get(argSetPrefixPositionPrefix),
				ctx.Bool("clear"),
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

			namespace, err = storage.FindNamespace(nsName)
			if err != nil {
				return resourceNotFoundError(nsName)
			}

			account, err = namespace.FindAccount(accName)
			if err != nil {
				return resourceNotFoundError(fmt.Sprintf("%s/%s", namespace.Name, accName))
			}

			account.Prefix = prefix

			return nil
		},
	}
}

func askForSetPrefixDetails(namespace, account, prefix string, isClear bool) (string, string, string) {
	term := terminal.New(os.Stdin, os.Stdout, os.Stderr)

	for len(namespace) < 1 {
		namespace, _ = term.Read("Namespace:")
	}

	for len(account) < 1 {
		account, _ = term.Read("Account:")
	}

	if isClear {
		return namespace, account, ""
	}

	for len(prefix) < 1 {
		prefix, _ = term.Read("Prefix:")
	}

	return namespace, account, prefix
}
