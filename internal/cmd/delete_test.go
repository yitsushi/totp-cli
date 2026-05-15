package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/yitsushi/totp-cli/internal/storage"
	"github.com/yitsushi/totp-cli/internal/terminal"
)

func TestDelete(t *testing.T) {
	suite.Run(t, &DeleteTestSuite{})
}

type DeleteTestSuite struct {
	suite.Suite
}

func (suite *DeleteTestSuite) TestDeleteAccountConfirmed() {
	ns := &storage.Namespace{
		Name:     "myns",
		Accounts: []*storage.Account{{Name: "github"}, {Name: "gitlab"}},
	}
	stor := newFakeStorage(ns)
	term := terminal.New(bytes.NewReader([]byte("yes\n")), &bytes.Buffer{}, &bytes.Buffer{})

	err := executeDelete(stor, term, "myns", "github")

	suite.Require().NoError(err)
	suite.True(stor.SaveCalled)
	suite.Len(ns.Accounts, 1)
	suite.Equal("gitlab", ns.Accounts[0].Name)
}

func (suite *DeleteTestSuite) TestDeleteAccountDenied() {
	ns := &storage.Namespace{
		Name:     "myns",
		Accounts: []*storage.Account{{Name: "github"}},
	}
	stor := newFakeStorage(ns)
	term := terminal.New(bytes.NewReader([]byte("no\n")), &bytes.Buffer{}, &bytes.Buffer{})

	err := executeDelete(stor, term, "myns", "github")

	suite.Require().NoError(err)
	suite.Len(ns.Accounts, 1, "account should not be removed when denied")
}

func (suite *DeleteTestSuite) TestDeleteNamespaceConfirmed() {
	stor := newFakeStorage(
		&storage.Namespace{Name: "myns"},
		&storage.Namespace{Name: "other"},
	)
	term := terminal.New(bytes.NewReader([]byte("yes\n")), &bytes.Buffer{}, &bytes.Buffer{})

	err := executeDelete(stor, term, "myns", "")

	suite.Require().NoError(err)
	suite.Len(stor.ListNamespaces(), 1)
	suite.Equal("other", stor.ListNamespaces()[0].Name)
}

func (suite *DeleteTestSuite) TestDeleteNamespaceNotFound() {
	stor := newFakeStorage()
	term := terminal.New(bytes.NewReader([]byte("yes\n")), &bytes.Buffer{}, &bytes.Buffer{})

	err := executeDelete(stor, term, "missing", "")

	suite.Require().ErrorIs(err, storage.NotFoundError{Type: "namespace", Name: "missing"})
}

func (suite *DeleteTestSuite) TestDeleteAccountNotFound() {
	ns := &storage.Namespace{Name: "myns"}
	stor := newFakeStorage(ns)
	term := terminal.New(bytes.NewReader([]byte("yes\n")), &bytes.Buffer{}, &bytes.Buffer{})

	err := executeDelete(stor, term, "myns", "missing")

	suite.Require().Error(err)
}
