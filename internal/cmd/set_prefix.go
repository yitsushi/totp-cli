package cmd

import (
	"fmt"
	"os"

	"github.com/yitsushi/go-commander"

	s "github.com/yitsushi/totp-cli/internal/storage"
	"github.com/yitsushi/totp-cli/internal/terminal"
)

// SetPrefix structure is the representation of the add-token command.
type SetPrefix struct{}

// Execute is the main function. It will be called on set-prefix command.
func (c *SetPrefix) Execute(opts *commander.CommandHelper) {
	var (
		namespace *s.Namespace
		account   *s.Account
		err       error
	)

	nsName, accName, prefix := c.askForSetPrefixDetails(opts)

	storage, err := s.PrepareStorage()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	defer func() {
		if err := storage.Save(); err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			os.Exit(1)
		}
	}()

	namespace, err = storage.FindNamespace(nsName)
	if err != nil {
		panic(fmt.Sprintf("%s does not exist!\n", nsName))
	}

	account, err = namespace.FindAccount(accName)
	if err != nil {
		panic(fmt.Sprintf("%s/%s does not exist!\n", namespace.Name, accName))
	}

	account.Prefix = prefix
}

// ArgumentDescription descripts the required and potential arguments.
func (c *SetPrefix) ArgumentDescription() string {
	return "[namespace] [account] [prefix]"
}

// Description will be displayed as Description (woooo) in the general help.
func (c *SetPrefix) Description() string {
	return "Set prefix for a token"
}

// Help is a general (human readable) command specific (long) help.
func (c *SetPrefix) Help() string {
	return ""
}

// Examples lists a few example as array. Will be used in the command specific help.
func (c *SetPrefix) Examples() []string {
	return []string{
		"",
		"mynamespace",
		"mynamespace myaccount",
		"mynamespace myaccount prefix",
	}
}

// NewSetPrefix creates a new SetPrefix command.
func NewSetPrefix(appName string) *commander.CommandWrapper {
	return &commander.CommandWrapper{
		Handler: &SetPrefix{},
		Help: &commander.CommandDescriptor{
			Name:             "set-prefix",
			ShortDescription: "Set prefix for a token",
			Arguments:        "[namespace] [account] [prefix]",
			Examples: []string{
				"",
				"mynamespace",
				"mynamespace myaccount",
				"mynamespace myaccount prefix",
			},
		},
	}
}

// Private functions

func (c *SetPrefix) askForSetPrefixDetails(opts *commander.CommandHelper) (namespace, account, prefix string) {
	term := terminal.New(os.Stdin, os.Stdout, os.Stderr)

	namespace = opts.Arg(0)
	for len(namespace) < 1 {
		namespace, _ = term.Read("Namespace:")
	}

	account = opts.Arg(1)
	for len(account) < 1 {
		account, _ = term.Read("Account:")
	}

	prefix = opts.Arg(2)
	for len(prefix) < 1 {
		prefix, _ = term.Read("Prefix:")
	}

	return
}
