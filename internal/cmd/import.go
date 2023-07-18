package cmd

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"

	s "github.com/yitsushi/totp-cli/internal/storage"
	"github.com/yitsushi/totp-cli/internal/terminal"
)

// ImportCommand is the import subcommand.
func ImportCommand() *cli.Command {
	return &cli.Command{
		Name:  "import",
		Usage: "Import tokens from a yaml file.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "input",
				Usage:    "Input YAML file. (REQUIRED)",
				Required: true,
			},
		},
		Action: func(ctx *cli.Context) (err error) {
			var (
				file    []byte
				storage *s.Storage
			)

			if file, err = os.ReadFile(ctx.String("input")); err != nil {
				return
			}

			nsList := []*s.Namespace{}

			if err = yaml.Unmarshal(file, &nsList); err != nil {
				err = CommandError{Message: "invalid file format"}

				return
			}

			if storage, err = s.PrepareStorage(); err != nil {
				return
			}

			defer func() {
				if err != nil {
					return
				}

				err = storage.Save()
			}()

			importNamespaces(storage, nsList)

			return nil
		},
	}
}

func importNamespaces(storage *s.Storage, nsList []*s.Namespace) {
	term := terminal.New(os.Stdin, os.Stdout, os.Stderr)

	for _, ns := range nsList {
		internalNS, err := storage.FindNamespace(ns.Name)
		if err != nil {
			storage.Namespaces = append(storage.Namespaces, ns)

			continue
		}

		for _, acc := range ns.Accounts {
			internalAcc, err := internalNS.FindAccount(acc.Name)
			if err != nil {
				internalNS.Accounts = append(internalNS.Accounts, acc)

				continue
			}

			msg := fmt.Sprintf(
				"[%s/%s] Account already exist, do you want to overwrite it?",
				ns.Name, acc.Name,
			)
			if term.Confirm(msg) {
				internalAcc.Token = acc.Token
				internalAcc.Prefix = acc.Prefix
			}
		}
	}
}
