package cmd

import (
	"fmt"

	"github.com/yitsushi/go-commander"

	s "github.com/yitsushi/totp-cli/internal/storage"
	"github.com/yitsushi/totp-cli/internal/util"
)

// AddToken structure is the representation of the add-token command.
type AddToken struct{}

// Execute is the main function. It will be called on add-token command.
func (c *AddToken) Execute(opts *commander.CommandHelper) {
	var (
		namespace *s.Namespace
		account   *s.Account
		err       error
	)

	nsName, accName, token := c.askForAddTokenDetails(opts)

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

// ArgumentDescription descripts the required and potential arguments.
func (c *AddToken) ArgumentDescription() string {
	return "[namespace] [account]"
}

// Description will be displayed as Description (woooo) in the general help.
func (c *AddToken) Description() string {
	return "Add new token"
}

// Help is a general (human readable) command specific (long) help.
func (c *AddToken) Help() string {
	return ""
}

// Examples lists a few example as array. Will be used in the command specific help.
func (c *AddToken) Examples() []string {
	return []string{
		"",
		"mynamespace",
		"mynamespace myaccount",
	}
}

// NewAddToken createa new AddToken command.
func NewAddToken(appName string) *commander.CommandWrapper {
	return &commander.CommandWrapper{
		Handler: &AddToken{},
		Help: &commander.CommandDescriptor{
			Name:             "add-token",
			ShortDescription: "Add new token",
			Arguments:        "[namespace] [account]",
			Examples: []string{
				"",
				"mynamespace",
				"mynamespace myaccount",
			},
		},
	}
}

// Private functions

func (c *AddToken) askForAddTokenDetails(opts *commander.CommandHelper) (namespace, account, token string) {
	namespace = opts.Arg(0)
	for len(namespace) < 1 {
		namespace = util.Ask("Namespace")
	}

	account = opts.Arg(1)
	for len(account) < 1 {
		account = util.Ask("Account")
	}

	for len(token) < 1 {
		token = util.Ask("Token")
	}

	return
}
