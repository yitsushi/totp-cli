package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/yitsushi/go-commander"

	s "github.com/yitsushi/totp-cli/internal/storage"
)

// Dump structure is the representation of the dump command.
type Dump struct{}

// Execute is the main function. It will be called on dump command.
func (c *Dump) Execute(opts *commander.CommandHelper) {
	storage, err := s.PrepareStorage()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	out, _ := json.Marshal(storage)

	fmt.Printf("%s\n", out)
}

// NewDump creates a new Dump command.
func NewDump(appName string) *commander.CommandWrapper {
	return &commander.CommandWrapper{
		Handler: &Dump{},
		Help: &commander.CommandDescriptor{
			Name:             "dump",
			ShortDescription: "Dump all available namespaces or accounts under a namespace",
			Arguments:        "[namespace]",
			Examples: []string{
				"",
			},
		},
	}
}
