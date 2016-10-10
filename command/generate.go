package command

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/Yitsushi/totp-cli/security"
	s "github.com/Yitsushi/totp-cli/storage"
	"github.com/Yitsushi/totp-cli/util"
)

// Generate structure is the representation of the generate command
type Generate struct {
}

// Description will be displayed as Description (woooo) in the general help
func (c *Generate) Description() string {
	return "Generate a specific OTP"
}

// ArgumentDescription descripts the required and potential arguments
func (c *Generate) ArgumentDescription() string {
	return "<namespace>.<account>"
}

// Execute is the main function. It will be called on generate command
func (c *Generate) Execute() {
	term := flag.Arg(1)
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

// Help is a general (human readable) command specific (long) help
func (c *Generate) Help() string {
	return ""
}

// Examples lists a few example as array. Will be used in the command specific help
func (c *Generate) Examples() []string {
	return []string{"mynamespace.myaccount"}
}
