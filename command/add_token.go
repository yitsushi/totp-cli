package command

import (
	"flag"
	"fmt"

	s "github.com/Yitsushi/totp-cli/storage"
	"github.com/Yitsushi/totp-cli/util"
)

// AddToken stucture is the representation of the add-token command
type AddToken struct {
}

// Execute is the main function. It will be called on add-token command
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

// ArgumentDescription descripts the required and potential arguments
func (c *AddToken) ArgumentDescription() string {
	return "[namespace] [account]"
}

// Description will be displayed as Description (woooo) in the general help
func (c *AddToken) Description() string {
	return fmt.Sprintf("Add new token")
}

// Help is a general (human readable) command specific (long) help
func (c *AddToken) Help() string {
	return ""
}

// Examples lists a few example as array. Will be used in the command specific help
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
