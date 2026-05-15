package cmd

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	s "github.com/yitsushi/totp-cli/internal/storage"
	"github.com/yitsushi/totp-cli/internal/terminal"
	"gopkg.in/yaml.v3"
)

// ImportCommand is the subcommand to import data from a YAML file.
func ImportCommand() *cli.Command {
	return &cli.Command{
		Name:  "import",
		Usage: "Import tokens from a yaml file.",
		Flags: []cli.Flag{
			flagInput(),
		},
		Action: func(ctx *cli.Context) error {
			file, err := os.ReadFile(ctx.String("input"))
			if err != nil {
				return fmt.Errorf("failed to read input file: %w", err)
			}

			var nsList []*s.Namespace

			err = yaml.Unmarshal(file, &nsList)
			if err != nil {
				return CommandError{Message: "invalid file format"}
			}

			term := terminal.New(os.Stdin, os.Stdout, os.Stderr)
			storage := prepareStorage(term)

			err = storage.Prepare()
			if err != nil {
				return err
			}

			executeImport(storage, term, nsList)

			return storage.Save()
		},
	}
}

func executeImport(storage s.Storage, term terminal.Terminal, nsList []*s.Namespace) {
	for _, ns := range nsList {
		internalNS, _ := storage.AddNamespace(&s.Namespace{Name: ns.Name})

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
