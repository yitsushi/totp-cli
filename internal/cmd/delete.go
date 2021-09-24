package cmd

import (
	"fmt"
	"os"

	"github.com/yitsushi/go-commander"

	s "github.com/yitsushi/totp-cli/internal/storage"
	"github.com/yitsushi/totp-cli/internal/util"
)

// Delete structure is the representation of the delete command.
type Delete struct{}

// Execute is the main function. It will be called on delete command.
func (c *Delete) Execute(opts *commander.CommandHelper) {
	var err error

	namespaceName := opts.Arg(0)
	if len(namespaceName) < 1 {
		panic("Wrong number of arguments")
	}

	accountName := opts.Arg(1)

	storage, err := s.PrepareStorage()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	defer func() {
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			os.Exit(1)
		}

		if err = storage.Save(); err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			os.Exit(1)
		}
	}()

	namespace, err := storage.FindNamespace(namespaceName)
	if err != nil {
		return
	}

	if accountName != "" {
		account, err := namespace.FindAccount(accountName)
		if err != nil {
			return
		}

		fmt.Printf("You want to delete '%s.%s' account.\n", namespace.Name, account.Name)

		if util.Confirm("Are you sure?") {
			namespace.DeleteAccount(account)

			return
		}
	} else {
		fmt.Printf("You want to delete '%s' namespace with %d accounts.\n", namespace.Name, len(namespace.Accounts))
		for _, account := range namespace.Accounts {
			fmt.Printf(" - %s.%s\n", namespace.Name, account.Name)
		}

		if util.Confirm("Are you sure?") {
			storage.DeleteNamespace(namespace)

			return
		}
	}
}

// NewDelete create a new Delete command.
func NewDelete(appName string) *commander.CommandWrapper {
	return &commander.CommandWrapper{
		Handler: &Delete{},
		Help: &commander.CommandDescriptor{
			Name:             "delete",
			ShortDescription: "Delete an account or a whole namespace",
			Arguments:        "<namespace> [account]",
			Examples: []string{
				"mynamespace",
				"mynamespace myaccount",
			},
		},
	}
}
