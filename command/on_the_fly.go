package command

import (
	"fmt"
	"time"

	"github.com/yitsushi/go-commander"
	"github.com/yitsushi/totp-cli/security"
	"github.com/yitsushi/totp-cli/util"
)

// OnTheFly structure is the representation of the generate command.
type OnTheFly struct {
}

// Execute is the main function. It will be called on generate command.
func (c *OnTheFly) Execute(opts *commander.CommandHelper) {
	token := util.Read()

	fmt.Println(security.GenerateOTPCode(token, time.Now()))
}

// NewOnTheFly creates a new OnTheFly command.
func NewOnTheFly(appName string) *commander.CommandWrapper {
	return &commander.CommandWrapper{
		Handler: &OnTheFly{},
		Help: &commander.CommandDescriptor{
			Name:             "on-the-fly",
			ShortDescription: "Generate OTP from token on the fly",
		},
	}
}
