package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/yitsushi/totp-cli/internal/storage"
	"github.com/yitsushi/totp-cli/internal/terminal"
)

func TestImport(t *testing.T) {
	suite.Run(t, &ImportTestSuite{})
}

type ImportTestSuite struct {
	suite.Suite
}

func (suite *ImportTestSuite) TestImportNewNamespaceNoPrompt() {
	stor := newFakeStorage()
	// No "yes\n" input needed — new namespaces should not prompt
	term := terminal.New(bytes.NewReader([]byte{}), &bytes.Buffer{}, &bytes.Buffer{})

	nsList := []*storage.Namespace{
		{Name: "newns", Accounts: []*storage.Account{
			{Name: "github", Token: "SECRET"},
		}},
	}

	executeImport(stor, term, nsList)

	ns, err := stor.FindNamespace("newns")
	suite.Require().NoError(err)
	suite.Len(ns.Accounts, 1)
	suite.Equal("SECRET", ns.Accounts[0].Token)
}

func (suite *ImportTestSuite) TestImportExistingAccountPromptsOverwrite() {
	existing := &storage.Namespace{
		Name:     "myns",
		Accounts: []*storage.Account{{Name: "github", Token: "OLD"}},
	}
	stor := newFakeStorage(existing)
	// User confirms overwrite
	term := terminal.New(bytes.NewReader([]byte("yes\n")), &bytes.Buffer{}, &bytes.Buffer{})

	nsList := []*storage.Namespace{
		{Name: "myns", Accounts: []*storage.Account{
			{Name: "github", Token: "NEW"},
		}},
	}

	executeImport(stor, term, nsList)

	acc, _ := existing.FindAccount("github")
	suite.Equal("NEW", acc.Token)
}

func (suite *ImportTestSuite) TestImportExistingAccountDeniedKeepsOld() {
	existing := &storage.Namespace{
		Name:     "myns",
		Accounts: []*storage.Account{{Name: "github", Token: "OLD"}},
	}
	stor := newFakeStorage(existing)
	// User denies overwrite
	term := terminal.New(bytes.NewReader([]byte("no\n")), &bytes.Buffer{}, &bytes.Buffer{})

	nsList := []*storage.Namespace{
		{Name: "myns", Accounts: []*storage.Account{
			{Name: "github", Token: "NEW"},
		}},
	}

	executeImport(stor, term, nsList)

	acc, _ := existing.FindAccount("github")
	suite.Equal("OLD", acc.Token)
}
