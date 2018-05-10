package command

import (
	"fmt"
	"strings"

	"github.com/Yitsushi/go-commander"
	s "github.com/Yitsushi/totp-cli/storage"
	"github.com/Yitsushi/totp-cli/util"
)

// Delete structure is the representation of the delete command
type Delete struct {
}

// Execute is the main function. It will be called on delete command
func (c *Delete) Execute(opts *commander.CommandHelper) {
	term := opts.Arg(0)
	if len(term) < 1 {
		panic("Wrong number of arguments")
	}

	path := strings.Split(term, ".")

	nsName := path[0]
	accName := ""

	if len(path) > 1 {
		accName = path[1]
	}

	storage := s.PrepareStorage()

	namespace, err := storage.FindNamespace(nsName)
	util.Check(err)

	if accName != "" {
		account, err := namespace.FindAccount(accName)
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

// NewDelete create a new Delete command
func NewDelete(appName string) *commander.CommandWrapper {
	return &commander.CommandWrapper{
		Handler: &Delete{},
		Help: &commander.CommandDescriptor{
			Name:             "delete",
			ShortDescription: "Delete an account or a whole namespace",
			Arguments:        "<namespace>[.account]",
			Examples: []string{
				"mynamespace",
				"mynamespace.maccount",
			},
		},
	}
}
