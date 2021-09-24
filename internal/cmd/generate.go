package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/yitsushi/go-commander"

	"github.com/yitsushi/totp-cli/internal/security"
	s "github.com/yitsushi/totp-cli/internal/storage"
)

// Generate structure is the representation of the generate command.
type Generate struct{}

// Execute is the main function. It will be called on generate command.
func (c *Generate) Execute(opts *commander.CommandHelper) {
	namespaceName := opts.Arg(0)
	if len(namespaceName) < 1 {
		panic("Namespace is not defined")
	}

	accountName := opts.Arg(1)
	if len(accountName) < 1 {
		panic("Account is not defined")
	}

	storage, err := s.PrepareStorage()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	namespace, err := storage.FindNamespace(namespaceName)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	account, err := namespace.FindAccount(accountName)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	code, err := security.GenerateOTPCode(account.Token, time.Now())
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}

	fmt.Println(code)
}

// NewGenerate creates a new Generate command.
func NewGenerate(appName string) *commander.CommandWrapper {
	return &commander.CommandWrapper{
		Handler: &Generate{},
		Help: &commander.CommandDescriptor{
			Name:             "generate",
			ShortDescription: "Generate a specific OTP",
			Arguments:        "<namespace> <account>",
			Examples:         []string{"mynamespace myaccount"},
		},
	}
}
