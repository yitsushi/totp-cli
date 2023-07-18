package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/yitsushi/totp-cli/internal/security"
	"github.com/yitsushi/totp-cli/internal/terminal"
)

// InstantCommand is the instant subcommand.
func InstantCommand() *cli.Command {
	return &cli.Command{
		Name:      "instant",
		Usage:     "Generate an OTP from TOTP_TOKEN or stdin without the Storage backend.",
		ArgsUsage: " ",
		Action: func(_ *cli.Context) error {
			token := os.Getenv("TOTP_TOKEN")
			if token == "" {
				term := terminal.New(os.Stdin, os.Stdout, os.Stderr)
				token, _ = term.Read("")
			}

			code, err := security.GenerateOTPCode(token, time.Now())
			if err != nil {
				return err
			}

			fmt.Println(code)

			return nil
		},
	}
}
