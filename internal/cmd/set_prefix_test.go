package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/yitsushi/totp-cli/internal/storage"
	"github.com/yitsushi/totp-cli/internal/terminal"
)

func TestSetPrefix(t *testing.T) {
	suite.Run(t, &SetPrefixTestSuite{})
}

type SetPrefixTestSuite struct {
	suite.Suite
}

func (suite *SetPrefixTestSuite) TestSetPrefix() {
	ns := &storage.Namespace{
		Name:     "myns",
		Accounts: []*storage.Account{{Name: "github"}},
	}
	stor := newFakeStorage(ns)

	err := executeSetPrefix(stor, "myns", "github", "GH:")

	suite.Require().NoError(err)
	suite.True(stor.SaveCalled)

	acc, _ := ns.FindAccount("github")
	suite.Equal("GH:", acc.Prefix)
}

func (suite *SetPrefixTestSuite) TestSetPrefixNamespaceNotFound() {
	stor := newFakeStorage()

	err := executeSetPrefix(stor, "missing", "github", "GH:")

	suite.Require().Error(err)
}

func (suite *SetPrefixTestSuite) TestSetPrefixAccountNotFound() {
	ns := &storage.Namespace{Name: "myns"}
	stor := newFakeStorage(ns)

	err := executeSetPrefix(stor, "myns", "missing", "GH:")

	suite.Require().Error(err)
}

func (suite *SetPrefixTestSuite) TestAskForSetPrefixDetailsEOFReturnsError() {
	term := terminal.New(bytes.NewReader([]byte{}), &bytes.Buffer{}, &bytes.Buffer{})

	_, _, _, err := askForSetPrefixDetails(term, "", "", "", false)

	suite.Require().Error(err)
}

func (suite *SetPrefixTestSuite) TestAskForSetPrefixDetailsReadsFromTerminal() {
	input := bytes.NewReader([]byte("myns\nmyacc\nGH:\n"))
	term := terminal.New(input, &bytes.Buffer{}, &bytes.Buffer{})

	ns, acc, prefix, err := askForSetPrefixDetails(term, "", "", "", false)

	suite.Require().NoError(err)
	suite.Equal("myns", ns)
	suite.Equal("myacc", acc)
	suite.Equal("GH:", prefix)
}

func (suite *SetPrefixTestSuite) TestAskForSetPrefixDetailsUsesProvidedArgs() {
	input := bytes.NewReader([]byte("GH:\n"))
	term := terminal.New(input, &bytes.Buffer{}, &bytes.Buffer{})

	ns, acc, prefix, err := askForSetPrefixDetails(term, "myns", "myacc", "", false)

	suite.Require().NoError(err)
	suite.Equal("myns", ns)
	suite.Equal("myacc", acc)
	suite.Equal("GH:", prefix)
}

func (suite *SetPrefixTestSuite) TestAskForSetPrefixDetailsClear() {
	term := terminal.New(bytes.NewReader([]byte{}), &bytes.Buffer{}, &bytes.Buffer{})

	ns, acc, prefix, err := askForSetPrefixDetails(term, "myns", "myacc", "", true)

	suite.Require().NoError(err)
	suite.Equal("myns", ns)
	suite.Equal("myacc", acc)
	suite.Empty(prefix)
}

func (suite *SetPrefixTestSuite) TestAskForSetPrefixDetailsEOFOnAccountReturnsError() {
	term := terminal.New(bytes.NewReader([]byte("myns\n")), &bytes.Buffer{}, &bytes.Buffer{})

	_, _, _, err := askForSetPrefixDetails(term, "", "", "", false)

	suite.Require().Error(err)
}

func (suite *SetPrefixTestSuite) TestAskForSetPrefixDetailsEOFOnPrefixReturnsError() {
	term := terminal.New(bytes.NewReader([]byte("myns\nmyacc\n")), &bytes.Buffer{}, &bytes.Buffer{})

	_, _, _, err := askForSetPrefixDetails(term, "", "", "", false)

	suite.Require().Error(err)
}
