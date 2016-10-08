package command

import (
	"flag"
	"fmt"

	s "github.com/Yitsushi/totp-cli/storage"
	"github.com/Yitsushi/totp-cli/util"
)

// List stucture is the representation of the list command
type List struct {
}

// Description will be displayed as Description (woooo) in the general help
func (c *List) Description() string {
	return "List all available namespaces or accounts under a namespace"
}

// ArgumentDescription descripts the required and potential arguments
func (c *List) ArgumentDescription() string {
	return "[namespace]"
}

// Execute is the main function. It will be called on list command
func (c *List) Execute() {
	storage := s.PrepareStorage()
	ns := flag.Arg(1)
	if len(ns) < 1 {
		for _, namespace := range storage.Namespaces {
			fmt.Printf("%s (Number of accounts: %d)\n", namespace.Name, len(namespace.Accounts))
		}

		return
	}

	namespace, err := storage.FindNamespace(ns)
	util.Check(err)

	for _, account := range namespace.Accounts {
		fmt.Printf("%s.%s\n", namespace.Name, account.Name)
	}
}

// Help is a general (human readable) command specific (long) help
func (c *List) Help() string {
	return ""
}

// Examples lists a few example as array. Will be used in the command specific help
func (c *List) Examples() []string {
	return []string{
		"",
		"mynamespace",
	}
}
