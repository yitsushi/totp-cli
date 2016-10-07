package command

import (
	"flag"
	"fmt"

	s "github.com/Yitsushi/totp-cli/storage"
	"github.com/Yitsushi/totp-cli/util"
)

type AddToken struct {
}

func (c *AddToken) Execute() {
	var namespace *s.Namespace
	var account *s.Account
	var err error

	nsName, accName, token := c.askForAddTokenDetails()

	storage := s.PrepareStorage()

	namespace, err = storage.FindNamespace(nsName)
	if err != nil {
		namespace = &s.Namespace{Name: nsName}
		storage.Namespaces = append(storage.Namespaces, namespace)
	}

	account, err = namespace.FindAccount(accName)
	if err == nil {
		fmt.Printf("%s.%s exists!\n", namespace.Name, account.Name)
	}

	account = &s.Account{Name: accName, Token: token}

	namespace.Accounts = append(namespace.Accounts, account)

	storage.Save()
}

func (c *AddToken) ArgumentDescription() string {
	return "[namespace] [account]"
}

func (c *AddToken) Description() string {
	return fmt.Sprintf("Add new token")
}

func (c *AddToken) Help() string {
	return ""
}

func (c *AddToken) Examples() []string {
	return []string{
		"",
		"mynamespace",
		"mynamespace myaccount",
	}
}

// Private functions

func (c *AddToken) askForAddTokenDetails() (namespace, account, token string) {
	namespace = flag.Arg(1)
	account = flag.Arg(2)
	for len(namespace) < 1 {
		namespace = util.Ask("Namespace")
	}
	for len(account) < 1 {
		account = util.Ask("Account")
	}
	for len(token) < 1 {
		token = util.Ask("Token")
	}

	return
}
