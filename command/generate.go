package command

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/Yitsushi/totp-cli/security"
	s "github.com/Yitsushi/totp-cli/storage"
)

type Generate struct {
}

func (c *Generate) Description() string {
	return "Generate a specific OTP"
}

func (c *Generate) ArgumentDescription() string {
	return "<namespace>.<account>"
}

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
	if err != nil {
		fmt.Println(err)
		return
	}
	account, err := namespace.FindAccount(path[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(security.GenerateOTPCode(account.Token, time.Now()))
}

func (c *Generate) Help() string {
	return ""
}

func (c *Generate) Examples() []string {
	return []string{"mynamespace.myaccount"}
}
