package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/yitsushi/totp-cli/internal/cmd"
	"github.com/yitsushi/totp-cli/internal/info"
)

func newApplication() *cli.App {
	stdErr := os.Stderr

	return &cli.App{
		Name:     info.AppName,
		HelpName: "totp-cli",
		Usage:    "Authy/Google Authenticator like TOTP CLI tool written in Go.",
		Version:  info.AppVersion,
		Commands: []*cli.Command{
			cmd.AddTokenCommand(),
			cmd.ChangePasswordCommand(),
			cmd.DeleteCommand(),
			cmd.DumpCommand(),
			cmd.GenerateCommand(),
			cmd.ImportCommand(),
			cmd.InstantCommand(),
			cmd.ListCommand(),
			cmd.SetPrefixCommand(),
			cmd.UpdateCommand(),
		},
		Flags: []cli.Flag{},
		Authors: []*cli.Author{
			{Name: "Efertone", Email: "efertone@pm.me"},
		},
		Copyright: "",
		ExitErrHandler: func(ctx *cli.Context, err error) {
			if err == nil {
				return
			}

			fmt.Fprintf(stdErr, " !!! %s\n", err)

			_ = cli.ShowAppHelp(ctx)
		},
		Suggest: true,
	}
}
