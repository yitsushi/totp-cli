package command

import (
	"fmt"
	"strings"
	"time"

	"github.com/Yitsushi/go-commander"
	"github.com/Yitsushi/totp-cli/security"
	s "github.com/Yitsushi/totp-cli/storage"
	"github.com/Yitsushi/totp-cli/util"
)

// Generate structure is the representation of the generate command
type Generate struct {
}

// Execute is the main function. It will be called on generate command
func (c *Generate) Execute(opts *commander.CommandHelper) {
	term := opts.Arg(0)
	if len(term) < 1 {
		panic("Namespace and Account are not defined")
	}

	path := strings.Split(term, ".")

	if len(path) < 2 {
		panic("Account is not defined")
	}

	storage := s.PrepareStorage()

	namespace, err := storage.FindNamespace(path[0])
	util.Check(err)

	account, err := namespace.FindAccount(path[1])
	util.Check(err)

	fmt.Println(security.GenerateOTPCode(account.Token, time.Now()))
}

func NewGenerate(appName string) *commander.CommandWrapper {
	return &commander.CommandWrapper{
		Handler: &Generate{},
		Help: &commander.CommandDescriptor{
			Name:             "generate",
			ShortDescription: "Generate a specific OTP",
			Arguments:        "<namespace>.<account>",
			Examples:         []string{"mynamespace.myaccount"},
		},
	}
}
