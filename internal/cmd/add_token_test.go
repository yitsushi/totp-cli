package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/yitsushi/totp-cli/internal/storage"
	"github.com/yitsushi/totp-cli/internal/terminal"
)

func TestAddToken(t *testing.T) {
	suite.Run(t, &AddTokenTestSuite{})
}

type AddTokenTestSuite struct {
	suite.Suite
}

func (suite *AddTokenTestSuite) TestAddNewToken() {
	stor := newFakeStorage(&storage.Namespace{Name: "myns"})

	err := executeAddToken(stor, "myns", "github", "SECRETTOKEN", AccountOptions{
		Length:    6,
		Algorithm: "sha1",
	})

	suite.Require().NoError(err)
	suite.True(stor.SaveCalled)

	ns, _ := stor.FindNamespace("myns")
	acc, err := ns.FindAccount("github")
	suite.Require().NoError(err)
	suite.Equal("SECRETTOKEN", acc.Token)
	suite.Equal(uint(6), acc.Length)
}

func (suite *AddTokenTestSuite) TestAddTokenCreatesNamespaceIfMissing() {
	stor := newFakeStorage()

	err := executeAddToken(stor, "newns", "github", "TOKEN", AccountOptions{})

	suite.Require().NoError(err)

	ns, err := stor.FindNamespace("newns")
	suite.Require().NoError(err)
	suite.Len(ns.Accounts, 1)
}

func (suite *AddTokenTestSuite) TestAddDuplicateTokenReturnsError() {
	ns := &storage.Namespace{
		Name:     "myns",
		Accounts: []*storage.Account{{Name: "github", Token: "OLD"}},
	}
	stor := newFakeStorage(ns)

	err := executeAddToken(stor, "myns", "github", "NEW", AccountOptions{})

	suite.Require().Error(err)
	suite.False(stor.SaveCalled)
}

func (suite *AddTokenTestSuite) TestAskForAddTokenDetailsEOFReturnsError() {
	term := terminal.New(bytes.NewReader([]byte{}), &bytes.Buffer{}, &bytes.Buffer{})

	_, _, _, err := askForAddTokenDetails(term, "", "")

	suite.Require().Error(err)
}

func (suite *AddTokenTestSuite) TestAskForAddTokenDetailsEOFOnAccountReturnsError() {
	term := terminal.New(bytes.NewReader([]byte("myns\n")), &bytes.Buffer{}, &bytes.Buffer{})

	_, _, _, err := askForAddTokenDetails(term, "", "")

	suite.Require().Error(err)
}

func (suite *AddTokenTestSuite) TestAskForAddTokenDetailsReadsFromTerminal() {
	input := bytes.NewReader([]byte("mynamespace\nmyaccount\nmytoken\n"))
	term := terminal.New(input, &bytes.Buffer{}, &bytes.Buffer{})

	ns, acc, token, err := askForAddTokenDetails(term, "", "")

	suite.Require().NoError(err)
	suite.Equal("mynamespace", ns)
	suite.Equal("myaccount", acc)
	suite.Equal("mytoken", token)
}

func (suite *AddTokenTestSuite) TestAskForAddTokenDetailsUsesProvidedArgs() {
	input := bytes.NewReader([]byte("mytoken\n"))
	term := terminal.New(input, &bytes.Buffer{}, &bytes.Buffer{})

	ns, acc, token, err := askForAddTokenDetails(term, "ns", "acc")

	suite.Require().NoError(err)
	suite.Equal("ns", ns)
	suite.Equal("acc", acc)
	suite.Equal("mytoken", token)
}
