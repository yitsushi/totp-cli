package cmd

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	s "github.com/yitsushi/totp-cli/internal/storage"
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
			&cli.BoolFlag{
				Name:  "yes-please",
				Value: false,
				Usage: warningMsg,
			},
			&cli.StringFlag{
				Name:     "output",
				Usage:    "Output file. (REQUIRED)",
				Required: true,
			},
		},
		Action: func(ctx *cli.Context) error {
			if !ctx.Bool("yes-please") {
				return CommandError{
					Message: warningMsg,
				}
			}

			storage := s.NewFileStorage()
			if err := storage.Prepare(); err != nil {
				return err
			}

			out, err := yaml.Marshal(storage.ListNamespaces())
			if err != nil {
				return fmt.Errorf("failed to marshal storage: %w", err)
			}

			if err := os.WriteFile(ctx.String("output"), out, strictDumpFilePerms); err != nil {
				return fmt.Errorf("failed to write output file: %w", err)
			}

			return nil
		},
	}
}
