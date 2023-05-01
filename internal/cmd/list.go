package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/yitsushi/go-commander"

	s "github.com/yitsushi/totp-cli/internal/storage"
)

// List structure is the representation of the list command.
type List struct{}

// Execute is the main function. It will be called on list command.
func (c *List) Execute(opts *commander.CommandHelper) {
	storage, err := s.PrepareStorage()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	ns := opts.Arg(0)
	if len(ns) < 1 {
		for _, namespace := range storage.Namespaces {
			fmt.Printf("%s (Number of accounts: %d)\n", namespace.Name, len(namespace.Accounts))
		}

		return
	}

	namespace, err := storage.FindNamespace(ns)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		os.Exit(1)
	}

	sort.Slice(namespace.Accounts, func(i, j int) bool {
		return namespace.Accounts[i].Name < namespace.Accounts[j].Name
	})

	for _, account := range namespace.Accounts {
		fmt.Printf("%s.%s\n", namespace.Name, account.Name)
	}
}

// NewList creates a new List command.
func NewList(appName string) *commander.CommandWrapper {
	return &commander.CommandWrapper{
		Handler: &List{},
		Help: &commander.CommandDescriptor{
			Name:             "list",
			ShortDescription: "List all available namespaces or accounts under a namespace",
			Arguments:        "[namespace]",
			Examples: []string{
				"",
				"mynamespace",
			},
		},
	}
}
