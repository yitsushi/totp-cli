package cmd

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	s "github.com/yitsushi/totp-cli/internal/storage"
	"github.com/yitsushi/totp-cli/internal/terminal"
	"gopkg.in/yaml.v3"
)

// DumpCommand is the dump subcommand.
func DumpCommand() *cli.Command {
	warningMsg := "The output is NOT encrypted. Use this flag to verify you really want to dump all secrets."

	return &cli.Command{
		Name:      "dump",
		Usage:     "Dump all available accounts under all namespaces.",
		ArgsUsage: " ",
		Flags: []cli.Flag{
			flagYesPlease(warningMsg),
			flagOutput(),
		},
		Action: func(ctx *cli.Context) error {
			if !ctx.Bool("yes-please") {
				return CommandError{Message: warningMsg}
			}

			term := terminal.New(os.Stdin, os.Stdout, os.Stderr)
			storage := prepareStorage(term)

			err := storage.Prepare()
			if err != nil {
				return err
			}

			return executeDump(storage, ctx.String("output"))
		},
	}
}

func executeDump(storage s.Storage, outputPath string) error {
	out, err := yaml.Marshal(storage.ListNamespaces())
	if err != nil {
		return fmt.Errorf("failed to marshal storage: %w", err)
	}

	err = os.WriteFile(outputPath, out, strictDumpFilePerms)
	if err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	return nil
}
