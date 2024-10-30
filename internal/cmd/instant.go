package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/yitsushi/totp-cli/internal/security"
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
			&cli.UintFlag{
				Name:  "length",
				Value: storage.DefaultTokenLength,
				Usage: "Length of the generated token.",
			},
			&cli.BoolFlag{
				Name:  "show-remaining",
				Value: false,
				Usage: "Show how much time left until the code will be invalid.",
			},
		},
		Action: func(ctx *cli.Context) error {
			token := os.Getenv("TOTP_TOKEN")
			if token == "" {
				term := terminal.New(os.Stdin, os.Stdout, os.Stderr)
				token, _ = term.Read("")
			}

			length := ctx.Uint("length")

			code, remaining, err := security.GenerateOTPCode(token, time.Now(), length)
			if err != nil {
				return err
			}

			fmt.Println(formatCode(code, remaining, ctx.Bool("show-remaining")))

			return nil
		},
	}
}
