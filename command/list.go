package command

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/ktr0731/go-fuzzyfinder"
	"github.com/yitsushi/go-commander"
	"github.com/yitsushi/totp-cli/security"
	s "github.com/yitsushi/totp-cli/storage"
	"github.com/yitsushi/totp-cli/util"
)

// List structure is the representation of the list command.
type List struct {
}

// Execute is the main function. It will be called on list command.
func (c *List) Execute(opts *commander.CommandHelper) {
	storage := s.PrepareStorage()

	var namespaceIdx int = -1
	var namespace *s.Namespace
	err := errors.New("")

	ns := opts.Arg(0)
	if len(ns) < 1 {

		namespaceIdx, err = fuzzyfinder.Find(storage.Namespaces, func(i int) string {
			return storage.Namespaces[i].Name
		})

		if err == fuzzyfinder.ErrAbort {
			fmt.Println("No Selection")
			return
		}

		util.Check(err)
	}

	if namespaceIdx == -1 {
		namespace, err = storage.FindNamespace(ns)
		util.Check(err)
	} else {
		namespace = storage.Namespaces[namespaceIdx]
	}

	var otps []*string
	accountIdx := -1

	accountIdx, err = fuzzyfinder.Find(namespace.Accounts,
		func(i int) string {
			now := time.Now()
			timer := uint64(math.Floor(float64(now.Unix()) / float64(30)))
			secondsUntilInvalid := ((timer+1)*30 - uint64(now.Unix()))

			account, _ := namespace.FindAccount(namespace.Accounts[i].Name)
			generatedOtp := string(security.GenerateOTPCode(account.Token, now))
			otps = append(otps, &generatedOtp)
			return namespace.Accounts[i].Name + strings.Repeat(" ", (10-len(namespace.Accounts[i].Name))) + "  |  " + generatedOtp + "  |  " + strconv.Itoa(int(secondsUntilInvalid))
		})

	if err == fuzzyfinder.ErrAbort {
		fmt.Println("No Selection")
		return
	}

	fmt.Println(*otps[accountIdx])

}

// NewList creates a new List command.
func NewList(appName string) *commander.CommandWrapper {
	return &commander.CommandWrapper{
		Handler: &List{},
		Help: &commander.CommandDescriptor{
			Name:             "list",
			ShortDescription: "List all available namespaces or accounts under a namespace",
			Arguments:        "[namespace]",
			Examples: []string{
				"",
				"mynamespace",
			},
		},
	}
}
