package command

import (
	"fmt"

	"github.com/yitsushi/go-commander"
	s "github.com/yitsushi/totp-cli/storage"
	"github.com/yitsushi/totp-cli/util"
)

// Delete structure is the representation of the delete command.
type Delete struct {
}

// Execute is the main function. It will be called on delete command.
func (c *Delete) Execute(opts *commander.CommandHelper) {
	namespaceName := opts.Arg(0)
	if len(namespaceName) < 1 {
		panic("Wrong number of arguments")
	}

	accountName := opts.Arg(1)
	storage := s.PrepareStorage()

	namespace, err := storage.FindNamespace(namespaceName)
	util.Check(err)

	if accountName != "" {
		account, err := namespace.FindAccount(accountName)
		util.Check(err)

		fmt.Printf("You want to delete '%s.%s' account.\n", namespace.Name, account.Name)

		if util.Confirm("Are you sure?") {
			namespace.DeleteAccount(account)
			storage.Save()
		}
	} else {
		fmt.Printf("You want to delete '%s' namespace with %d accounts.\n", namespace.Name, len(namespace.Accounts))
		for _, account := range namespace.Accounts {
			fmt.Printf(" - %s.%s\n", namespace.Name, account.Name)
		}

		if util.Confirm("Are you sure?") {
			storage.DeleteNamespace(namespace)
			storage.Save()
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
