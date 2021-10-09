package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/yitsushi/go-commander"

	"github.com/yitsushi/totp-cli/internal/security"
	"github.com/yitsushi/totp-cli/internal/terminal"
)

// Instant structure is the representation of the instant command.
type Instant struct{}

// Execute is the main function. It will be called on instant command.
func (c *Instant) Execute(opts *commander.CommandHelper) {
	token := os.Getenv("TOTP_TOKEN")
	if token == "" {
		term := terminal.New(os.Stdin, os.Stdout, os.Stderr)
		token, _ = term.Read("")
	}

	code, err := security.GenerateOTPCode(token, time.Now())
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Println(code)
}

// NewInstant creates a new Instant command.
func NewInstant(appName string) *commander.CommandWrapper {
	return &commander.CommandWrapper{
		Handler: &Instant{},
		Help: &commander.CommandDescriptor{
			Name:             "instant",
			ShortDescription: "Generate an OTP from TOTP_TOKEN or stdin without the Storage backend",
		},
	}
}
