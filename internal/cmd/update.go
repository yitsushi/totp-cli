package cmd

import (
	"fmt"
	"slices"

	"github.com/urfave/cli/v2"
	s "github.com/yitsushi/totp-cli/internal/storage"
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
				return CommandError{Message: "namespace is not defined"}
			}

			accountName := ctx.Args().Get(argSetPrefixPositionAccount)
			if len(accountName) < 1 {
				return CommandError{Message: "account is not defined"}
			}

			storage := s.NewFileStorage()
			if err := storage.Prepare(); err != nil {
				return err
			}

			account, err := getAccount(storage, namespaceName, accountName)
			if err != nil {
				return err
			}

			if slices.Contains(ctx.LocalFlagNames(), "algorithm") {
				account.Algorithm = ctx.String("algorithm")
			}

			if slices.Contains(ctx.LocalFlagNames(), "length") {
				account.Length = ctx.Uint("length")
			}

			if slices.Contains(ctx.LocalFlagNames(), "time-period") {
				account.TimePeriod = ctx.Int64("time-period")
			}

			if slices.Contains(ctx.LocalFlagNames(), "prefix") {
				account.Prefix = ctx.String("prefix")
			}

			if err := storage.Save(); err != nil {
				return fmt.Errorf("failed to save the storage: %w", err)
			}

			return nil
		},
	}
}
