package command

import (
	"fmt"
	"os"
	"time"

	"github.com/yitsushi/go-commander"
	"github.com/yitsushi/totp-cli/security"
)

// Instant structure is the representation of the instant command.
type Instant struct {
}

// Execute is the main function. It will be called on instant command.
func (c *Instant) Execute(opts *commander.CommandHelper) {
	token := os.Getenv("TOTP_TOKEN")
	if token == "" {
		panic("TOTP_TOKEN not defined")
	}

	fmt.Println(security.GenerateOTPCode(token, time.Now()))
}

// NewInstant creates a new Instant command.
func NewInstant(appName string) *commander.CommandWrapper {
	return &commander.CommandWrapper{
		Handler: &Instant{},
		Help: &commander.CommandDescriptor{
			Name:             "instant",
			ShortDescription: "Generate an OTP from TOTP_TOKEN without the Storage backend",
		},
	}
}
