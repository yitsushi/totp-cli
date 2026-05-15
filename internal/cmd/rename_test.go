package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/yitsushi/totp-cli/internal/storage"
	"github.com/yitsushi/totp-cli/internal/terminal"
)

func TestRename(t *testing.T) {
	suite.Run(t, &RenameTestSuite{})
}

type RenameTestSuite struct {
	suite.Suite
}

func (suite *RenameTestSuite) TestRenameNamespace() {
	stor := newFakeStorage(&storage.Namespace{Name: "old"})

	err := executeRenameNamespace(stor, "old", "new")

	suite.Require().NoError(err)
	suite.True(stor.SaveCalled)

	_, err = stor.FindNamespace("new")
	suite.Require().NoError(err)
}

func (suite *RenameTestSuite) TestRenameNamespaceNotFound() {
	stor := newFakeStorage()

	err := executeRenameNamespace(stor, "missing", "new")

	suite.Require().Error(err)
}

func (suite *RenameTestSuite) TestRenameAccount() {
	ns := &storage.Namespace{
		Name:     "myns",
		Accounts: []*storage.Account{{Name: "old-acc"}},
	}
	stor := newFakeStorage(ns)

	err := executeRenameAccount(stor, "myns", "old-acc", "new-acc")

	suite.Require().NoError(err)
	suite.True(stor.SaveCalled)

	_, err = ns.FindAccount("new-acc")
	suite.Require().NoError(err)
}

func (suite *RenameTestSuite) TestRenameAccountNamespaceNotFound() {
	stor := newFakeStorage()

	err := executeRenameAccount(stor, "missing", "acc", "new")

	suite.Require().Error(err)
}

func (suite *RenameTestSuite) TestAskForNamespaceRenameDetailsEOFReturnsError() {
	term := terminal.New(bytes.NewReader([]byte{}), &bytes.Buffer{}, &bytes.Buffer{})

	_, _, err := askForNamespaceRenameDetails(term, "", "")

	suite.Require().Error(err)
}

func (suite *RenameTestSuite) TestAskForAccountRenameDetailsEOFReturnsError() {
	term := terminal.New(bytes.NewReader([]byte{}), &bytes.Buffer{}, &bytes.Buffer{})

	_, _, _, err := askForAccountRenameDetails(term, "", "", "")

	suite.Require().Error(err)
}

func (suite *RenameTestSuite) TestRenameAccountNotFound() {
	ns := &storage.Namespace{Name: "myns"}
	stor := newFakeStorage(ns)

	err := executeRenameAccount(stor, "myns", "missing", "new-acc")

	suite.Require().Error(err)
}

func (suite *RenameTestSuite) TestAskForNamespaceRenameDetailsReadsFromTerminal() {
	input := bytes.NewReader([]byte("myns\nnewname\n"))
	term := terminal.New(input, &bytes.Buffer{}, &bytes.Buffer{})

	ns, newName, err := askForNamespaceRenameDetails(term, "", "")

	suite.Require().NoError(err)
	suite.Equal("myns", ns)
	suite.Equal("newname", newName)
}

func (suite *RenameTestSuite) TestAskForNamespaceRenameDetailsEOFOnNewNameReturnsError() {
	term := terminal.New(bytes.NewReader([]byte("myns\n")), &bytes.Buffer{}, &bytes.Buffer{})

	_, _, err := askForNamespaceRenameDetails(term, "", "")

	suite.Require().Error(err)
}

func (suite *RenameTestSuite) TestAskForAccountRenameDetailsEOFOnAccountReturnsError() {
	term := terminal.New(bytes.NewReader([]byte("myns\n")), &bytes.Buffer{}, &bytes.Buffer{})

	_, _, _, err := askForAccountRenameDetails(term, "", "", "")

	suite.Require().Error(err)
}

func (suite *RenameTestSuite) TestAskForAccountRenameDetailsEOFOnNewNameReturnsError() {
	term := terminal.New(bytes.NewReader([]byte("myns\nmyacc\n")), &bytes.Buffer{}, &bytes.Buffer{})

	_, _, _, err := askForAccountRenameDetails(term, "", "", "")

	suite.Require().Error(err)
}

func (suite *RenameTestSuite) TestAskForAccountRenameDetailsReadsFromTerminal() {
	input := bytes.NewReader([]byte("myns\nmyacc\nnewname\n"))
	term := terminal.New(input, &bytes.Buffer{}, &bytes.Buffer{})

	ns, acc, newName, err := askForAccountRenameDetails(term, "", "", "")

	suite.Require().NoError(err)
	suite.Equal("myns", ns)
	suite.Equal("myacc", acc)
	suite.Equal("newname", newName)
}
