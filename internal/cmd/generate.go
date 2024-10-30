package cmd

import (
	"fmt"
	"time"

	"github.com/urfave/cli/v2"
	s "github.com/yitsushi/totp-cli/internal/storage"
)

// GenerateCommand is the subcommand to generate a TOTP token.
func GenerateCommand() *cli.Command {
	return &cli.Command{
		Name:    "generate",
		Aliases: []string{"g"},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "follow",
				Value: false,
				Usage: "Generate codes continuously.",
			},
			&cli.BoolFlag{
				Name:  "show-remaining",
				Value: false,
				Usage: "Show how much time left until the code will be invalid.",
			},
		},
		Usage:     "Generate a specific OTP",
		ArgsUsage: "<namespace> <account>",
		Action: func(ctx *cli.Context) error {
			namespaceName := ctx.Args().Get(argSetPrefixPositionNamespace)
			if len(namespaceName) < 1 {
				return CommandError{Message: "namespace is not defined"}
			}

			accountName := ctx.Args().Get(argSetPrefixPositionAccount)
			if len(accountName) < 1 {
				return CommandError{Message: "account is not defined"}
			}

			follow := ctx.Bool("follow")

			storage := s.NewFileStorage()
			if err := storage.Prepare(); err != nil {
				return err
			}

			account, err := getAccount(storage, namespaceName, accountName)
			if err != nil {
				return err
			}

			code, remaining := generateCode(account)
			fmt.Println(formatCode(code, remaining, ctx.Bool("show-remaining")))

			if !follow {
				return nil
			}

			previousCode := code

			for {
				code, remaining := generateCode(account)
				if code != previousCode {
					fmt.Println(formatCode(code, remaining, ctx.Bool("show-remaining")))
					previousCode = code
				}

				time.Sleep(time.Second)
			}
		},
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
