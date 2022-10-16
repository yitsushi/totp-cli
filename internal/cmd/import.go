package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/yitsushi/go-commander"
	"gopkg.in/yaml.v3"

	s "github.com/yitsushi/totp-cli/internal/storage"
	"github.com/yitsushi/totp-cli/internal/terminal"
)

// Import structure is the representation of the import command.
type Import struct{}

// Execute is the main function. It will be called on instant command.
func (c *Import) Execute(opts *commander.CommandHelper) {
	if opts.Arg(0) == "" {
		opts.Log(DownloadError{Message: "wrong number of argument"}.Error())

		return
	}

	file, err := ioutil.ReadFile(opts.Arg(0))
	if err != nil {
		opts.Log(DownloadError{Message: "failed to read file"}.Error())

		return
	}

	nsList := []*s.Namespace{}

	err = yaml.Unmarshal(file, &nsList)
	if err != nil {
		opts.Log(DownloadError{Message: "invalid file format"}.Error())

		return
	}

	storage, err := s.PrepareStorage()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	defer func() {
		if err := storage.Save(); err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			os.Exit(1)
		}
	}()

	c.importNamespaces(storage, nsList)
}

func (c *Import) importNamespaces(storage *s.Storage, nsList []*s.Namespace) {
	term := terminal.New(os.Stdin, os.Stdout, os.Stderr)

	for _, ns := range nsList {
		internalNS, err := storage.FindNamespace(ns.Name)
		if err != nil {
			storage.Namespaces = append(storage.Namespaces, ns)

			continue
		}

		for _, acc := range ns.Accounts {
			internalAcc, err := internalNS.FindAccount(acc.Name)
			if err != nil {
				internalNS.Accounts = append(internalNS.Accounts, acc)

				continue
			}

			msg := fmt.Sprintf(
				"[%s/%s] Account already exist, do you want to overwrite it?",
				ns.Name, acc.Name,
			)
			if term.Confirm(msg) {
				internalAcc.Token = acc.Token
			}
		}
	}
}

// NewImport creates a new Instant command.
func NewImport(appName string) *commander.CommandWrapper {
	return &commander.CommandWrapper{
		Handler: &Import{},
		Help: &commander.CommandDescriptor{
			Name:             "import",
			ShortDescription: "Import tokens from a yaml file.",
			Arguments:        "<input-file>",
			Examples: []string{
				"credentials.yaml",
			},
		},
	}
}
