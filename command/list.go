package command

import (
	"flag"
	"fmt"

	s "github.com/Yitsushi/totp-cli/storage"
	"github.com/Yitsushi/totp-cli/util"
)

type List struct {
}

func (c *List) Description() string {
	return "List all available namespaces or accounts under a namespace"
}

func (c *List) ArgumentDescription() string {
	return "[namespace]"
}

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

func (c *List) Help() string {
	return ""
}

func (c *List) Examples() []string {
	return []string{
		"",
		"mynamespace",
	}
}
