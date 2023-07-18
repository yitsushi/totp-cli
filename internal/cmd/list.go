package cmd

import (
	"fmt"
	"sort"

	"github.com/urfave/cli/v2"

	s "github.com/yitsushi/totp-cli/internal/storage"
)

func ListCommand() *cli.Command {
	return &cli.Command{
		Name:      "list",
		Usage:     "List all available namespaces or accounts under a namespace.",
		ArgsUsage: "[namespace]",
		Action: func(ctx *cli.Context) error {
			storage, err := s.PrepareStorage()
			if err != nil {
				return nil
			}

			ns := ctx.Args().Get(argPositionNamespace)
			if len(ns) < 1 {
				for _, namespace := range storage.Namespaces {
					fmt.Printf("%s (Number of accounts: %d)\n", namespace.Name, len(namespace.Accounts))
				}

				return nil
			}

			namespace, err := storage.FindNamespace(ns)
			if err != nil {
				return err
			}

			sort.Slice(namespace.Accounts, func(i, j int) bool {
				return namespace.Accounts[i].Name < namespace.Accounts[j].Name
			})

			for _, account := range namespace.Accounts {
				fmt.Printf("%s.%s\n", namespace.Name, account.Name)
			}

			return nil
		},
	}
}
