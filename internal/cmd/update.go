package cmd

import (
	"fmt"
	"os"
	"slices"

	"github.com/urfave/cli/v2"
	s "github.com/yitsushi/totp-cli/internal/storage"
	"github.com/yitsushi/totp-cli/internal/terminal"
)

// UpdateCommand is the subcommand to update a TOTP token options.
func UpdateCommand() *cli.Command {
	return &cli.Command{
		Name:  "update",
		Usage: "Update the TOTP token options.",
		Flags: []cli.Flag{
			flagLength(),
			flagPrefix(),
			flagAlgorithm(),
			flagTimePeriod(),
		},
		ArgsUsage: "[namespace] [account]",
		Action: func(ctx *cli.Context) error {
			namespaceName := ctx.Args().Get(argSetPrefixPositionNamespace)
			if len(namespaceName) < 1 {
				return CommandError{Message: errMsgNamespaceNotDefined}
			}

			accountName := ctx.Args().Get(argSetPrefixPositionAccount)
			if len(accountName) < 1 {
				return CommandError{Message: "account is not defined"}
			}

			term := terminal.New(os.Stdin, os.Stdout, os.Stderr)
			storage := prepareStorage(term)

			err := storage.Prepare()
			if err != nil {
				return err
			}

			opts := AccountOptions{}

			if slices.Contains(ctx.LocalFlagNames(), "algorithm") {
				opts.Algorithm = ctx.String("algorithm")
			}

			if slices.Contains(ctx.LocalFlagNames(), "length") {
				opts.Length = ctx.Uint("length")
			}

			if slices.Contains(ctx.LocalFlagNames(), "time-period") {
				opts.TimePeriod = ctx.Int64("time-period")
			}

			if slices.Contains(ctx.LocalFlagNames(), "prefix") {
				opts.Prefix = ctx.String("prefix")
			}

			return executeUpdate(storage, namespaceName, accountName, opts, ctx.LocalFlagNames())
		},
	}
}

func executeUpdate(storage s.Storage, nsName, accName string, opts AccountOptions, setFlags []string) error {
	account, err := getAccount(storage, nsName, accName)
	if err != nil {
		return err
	}

	if slices.Contains(setFlags, "algorithm") {
		account.Algorithm = opts.Algorithm
	}

	if slices.Contains(setFlags, "length") {
		account.Length = opts.Length
	}

	if slices.Contains(setFlags, "time-period") {
		account.TimePeriod = opts.TimePeriod
	}

	if slices.Contains(setFlags, "prefix") {
		account.Prefix = opts.Prefix
	}

	err = storage.Save()
	if err != nil {
		return fmt.Errorf("failed to save the storage: %w", err)
	}

	return nil
}
