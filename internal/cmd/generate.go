package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli/v2"
	s "github.com/yitsushi/totp-cli/internal/storage"
	"github.com/yitsushi/totp-cli/internal/terminal"
)

// GenerateCommand is the subcommand to generate a TOTP token.
func GenerateCommand() *cli.Command {
	return &cli.Command{
		Name:    "generate",
		Aliases: []string{"g"},
		Flags: []cli.Flag{
			flagFollow(),
			flagShowRemaining(),
		},
		Usage:     "Generate a specific OTP",
		ArgsUsage: "<namespace> <account>",
		Action: func(ctx *cli.Context) error {
			namespaceName := ctx.Args().Get(argSetPrefixPositionNamespace)
			if len(namespaceName) < 1 {
				return CommandError{Message: errMsgNamespaceNotDefined}
			}

			accountName := ctx.Args().Get(argSetPrefixPositionAccount)
			if len(accountName) < 1 {
				return CommandError{Message: "account is not defined"}
			}

			term := terminal.New(os.Stdin, os.Stdout, os.Stderr)
			storage := prepareStorage(term)

			err := storage.Prepare()
			if err != nil {
				return err
			}

			return executeGenerate(storage, namespaceName, accountName, ctx.Bool("follow"), ctx.Bool("show-remaining"))
		},
	}
}

func executeGenerate(storage s.Storage, nsName, accName string, follow, showRemaining bool) error {
	account, err := getAccount(storage, nsName, accName)
	if err != nil {
		return err
	}

	code, remaining, err := generateCode(account)
	if err != nil {
		return err
	}

	fmt.Println(formatCode(code, remaining, showRemaining))

	if !follow {
		return nil
	}

	previousCode := code

	for {
		code, remaining, err := generateCode(account)
		if err != nil {
			return err
		}

		if code != previousCode {
			fmt.Println(formatCode(code, remaining, showRemaining))
			previousCode = code
		}

		time.Sleep(time.Second)
	}
}

func getAccount(storage s.Storage, namespaceName, accountName string) (*s.Account, error) {
	namespace, err := storage.FindNamespace(namespaceName)
	if err != nil {
		return nil, err
	}

	account, err := namespace.FindAccount(accountName)
	if err != nil {
		return nil, err
	}

	return account, nil
}
