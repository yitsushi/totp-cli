package cmd

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/yitsushi/totp-cli/internal/storage"
)

func TestList(t *testing.T) {
	suite.Run(t, &ListTestSuite{})
}

type ListTestSuite struct {
	suite.Suite
}

func (suite *ListTestSuite) TestListAllNamespaces() {
	stor := newFakeStorage(
		&storage.Namespace{Name: "ns1"},
		&storage.Namespace{Name: "ns2"},
	)

	err := executeList(stor, "")

	suite.Require().NoError(err)
}

func (suite *ListTestSuite) TestListAccountsInNamespace() {
	stor := newFakeStorage(
		&storage.Namespace{Name: "ns1", Accounts: []*storage.Account{
			{Name: "b-account"},
			{Name: "a-account"},
		}},
	)

	err := executeList(stor, "ns1")

	suite.Require().NoError(err)
}

func (suite *ListTestSuite) TestListNamespaceNotFound() {
	stor := newFakeStorage()

	err := executeList(stor, "missing")

	suite.Require().ErrorIs(err, storage.NotFoundError{Type: "namespace", Name: "missing"})
}
