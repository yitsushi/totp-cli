package command

import (
	"flag"
	"fmt"
	"strings"

	s "github.com/Yitsushi/totp-cli/storage"
	"github.com/Yitsushi/totp-cli/util"
)

// Delete structure is the representation of the delete command
type Delete struct {
}

// Execute is the main function. It will be called on delete command
func (c *Delete) Execute() {
	term := flag.Arg(1)
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

// ArgumentDescription descripts the required and potential arguments
func (c *Delete) ArgumentDescription() string {
	return "<namespace>[.account]"
}

// Description will be displayed as Description (woooo) in the general help
func (c *Delete) Description() string {
	return "Delete an account or a whole namespace"
}

// Help is a general (human readable) command specific (long) help
func (c *Delete) Help() string {
	return ""
}

// Examples lists a few example as array. Will be used in the command specific help
func (c *Delete) Examples() []string {
	return []string{
		"mynamespace myaccount",
		"mynamespace",
	}
}
