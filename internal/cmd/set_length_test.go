package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/yitsushi/totp-cli/internal/storage"
	"github.com/yitsushi/totp-cli/internal/terminal"
)

func TestSetLength(t *testing.T) {
	suite.Run(t, &SetLengthTestSuite{})
}

type SetLengthTestSuite struct {
	suite.Suite
}

func (suite *SetLengthTestSuite) TestSetLength() {
	ns := &storage.Namespace{
		Name:     "myns",
		Accounts: []*storage.Account{{Name: "github"}},
	}
	stor := newFakeStorage(ns)

	err := executeSetLength(stor, "myns", "github", 8)

	suite.Require().NoError(err)
	suite.True(stor.SaveCalled)

	acc, _ := ns.FindAccount("github")
	suite.Equal(uint(8), acc.Length)
}

func (suite *SetLengthTestSuite) TestSetLengthNamespaceNotFound() {
	stor := newFakeStorage()

	err := executeSetLength(stor, "missing", "github", 8)

	suite.Require().Error(err)
}

func (suite *SetLengthTestSuite) TestSetLengthAccountNotFound() {
	ns := &storage.Namespace{Name: "myns"}
	stor := newFakeStorage(ns)

	err := executeSetLength(stor, "myns", "missing", 8)

	suite.Require().Error(err)
}

func (suite *SetLengthTestSuite) TestAskForSetLengthDetailsEOFReturnsError() {
	term := terminal.New(bytes.NewReader([]byte{}), &bytes.Buffer{}, &bytes.Buffer{})

	_, _, _, err := askForSetLengthDetails(term, "", "", "")

	suite.Require().Error(err)
}

func (suite *SetLengthTestSuite) TestAskForSetLengthDetailsReadsFromTerminal() {
	input := bytes.NewReader([]byte("myns\nmyacc\n8\n"))
	term := terminal.New(input, &bytes.Buffer{}, &bytes.Buffer{})

	ns, acc, length, err := askForSetLengthDetails(term, "", "", "")

	suite.Require().NoError(err)
	suite.Equal("myns", ns)
	suite.Equal("myacc", acc)
	suite.Equal(uint(8), length)
}

func (suite *SetLengthTestSuite) TestAskForSetLengthDetailsUsesProvidedArgs() {
	input := bytes.NewReader([]byte("8\n"))
	term := terminal.New(input, &bytes.Buffer{}, &bytes.Buffer{})

	ns, acc, length, err := askForSetLengthDetails(term, "myns", "myacc", "")

	suite.Require().NoError(err)
	suite.Equal("myns", ns)
	suite.Equal("myacc", acc)
	suite.Equal(uint(8), length)
}

func (suite *SetLengthTestSuite) TestAskForSetLengthDetailsRejectsInvalidThenAccepts() {
	input := bytes.NewReader([]byte("abc\n8\n"))
	term := terminal.New(input, &bytes.Buffer{}, &bytes.Buffer{})

	_, _, length, err := askForSetLengthDetails(term, "myns", "myacc", "")

	suite.Require().NoError(err)
	suite.Equal(uint(8), length)
}

func (suite *SetLengthTestSuite) TestAskForSetLengthDetailsEOFOnAccountReturnsError() {
	term := terminal.New(bytes.NewReader([]byte("myns\n")), &bytes.Buffer{}, &bytes.Buffer{})

	_, _, _, err := askForSetLengthDetails(term, "", "", "")

	suite.Require().Error(err)
}

func (suite *SetLengthTestSuite) TestAskForSetLengthDetailsEOFOnLengthReturnsError() {
	term := terminal.New(bytes.NewReader([]byte("myns\nmyacc\n")), &bytes.Buffer{}, &bytes.Buffer{})

	_, _, _, err := askForSetLengthDetails(term, "", "", "")

	suite.Require().Error(err)
}
