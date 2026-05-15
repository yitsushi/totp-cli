package cmd

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/yitsushi/totp-cli/internal/storage"
	"github.com/yitsushi/totp-cli/internal/terminal"
)

// InstantCommand is the instant subcommand.
func InstantCommand() *cli.Command {
	return &cli.Command{
		Name:      "instant",
		Usage:     "Generate an OTP from TOTP_TOKEN or stdin without the Storage backend.",
		ArgsUsage: " ",
		Flags: []cli.Flag{
			flagLength(),
			flagShowRemaining(),
			flagAlgorithm(),
			flagTimePeriod(),
		},
		Action: func(ctx *cli.Context) error {
			account := storage.Account{
				Name:       "instant",
				Token:      os.Getenv("TOTP_TOKEN"),
				Length:     ctx.Uint("length"),
				Algorithm:  ctx.String("algorithm"),
				TimePeriod: ctx.Int64("time-period"),
			}

			if account.Token == "" {
				term := terminal.New(os.Stdin, os.Stdout, os.Stderr)

				var err error

				account.Token, err = term.Read("")
				if err != nil {
					return err
				}
			}

			code, remaining, err := generateCode(&account)
			if err != nil {
				return err
			}

			fmt.Println(formatCode(code, remaining, ctx.Bool("show-remaining")))

			return nil
		},
	}
}
