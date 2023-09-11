package cmd

import (
	"fmt"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/yitsushi/totp-cli/internal/security"
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

			storage, err := s.PrepareStorage()
			if err != nil {
				return err
			}

			namespace, err := storage.FindNamespace(namespaceName)
			if err != nil {
				return err
			}

			account, err := namespace.FindAccount(accountName)
			if err != nil {
				return err
			}

			code := generateCode(account)
			fmt.Println(code)

			if !follow {
				return nil
			}

			previousCode := code

			for {
				code := generateCode(account)
				if code != previousCode {
					fmt.Println(code)
					previousCode = code
				}

				time.Sleep(time.Second)
			}
		},
	}
}

func generateCode(account *s.Account) string {
	code, err := security.GenerateOTPCode(account.Token, time.Now(), account.Length)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}

	if account.Prefix != "" {
		code = account.Prefix + code
	}

	return code
}
