package cmd

import (
	s "github.com/yitsushi/totp-cli/internal/storage"
	"github.com/yitsushi/totp-cli/internal/terminal"
)

// AccountOptions groups the optional per-account fields shared by add-token
// and update commands.
type AccountOptions struct {
	Prefix     string
	Algorithm  string
	TimePeriod int64
	Length     uint
}

// terminalPasswordProvider adapts a Terminal into a storage.PasswordProvider.
type terminalPasswordProvider struct {
	term terminal.Terminal
}

func (p terminalPasswordProvider) GetPassword(prompt string) (string, error) {
	return p.term.Hidden(prompt)
}

// prepareStorage creates a FileBackend wired to the given terminal for
// password prompts. The caller is responsible for calling Prepare().
func prepareStorage(term terminal.Terminal) *s.FileBackend {
	return s.NewFileStorage(
		s.WithPasswordProvider(terminalPasswordProvider{term: term}),
	)
}
