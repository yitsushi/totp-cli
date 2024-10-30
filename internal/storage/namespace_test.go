package storage_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/yitsushi/totp-cli/internal/storage"
)

func TestAccount(t *testing.T) {
	suite.Run(t, &AccountTestSuite{})
}

type AccountTestSuite struct {
	suite.Suite
}

func (suite *AccountTestSuite) TestFindAccount() {
	namespace := &storage.Namespace{
		Name: "myNamespace",
		Accounts: []*storage.Account{
			{Name: "Account1", Token: "token1"},
			{Name: "Account2", Token: "token2"},
			{Name: "Account3", Token: "token3"},
		},
	}

	account, err := namespace.FindAccount("Account1")

	suite.Require().NoError(err)
	suite.Equal("Account1", account.Name, "Found account name should be Account1")
}

func (suite *AccountTestSuite) TestAccountNotFound() {
	namespace := &storage.Namespace{
		Name: "myNamespace",
		Accounts: []*storage.Account{
			{Name: "Account1", Token: "token1"},
			{Name: "Account2", Token: "token2"},
			{Name: "Account3", Token: "token3"},
		},
	}

	account, err := namespace.FindAccount("AccountNotFound")

	suite.Require().ErrorIs(err, storage.NotFoundError{Type: "account", Name: "AccountNotFound"})
	suite.Nil(account)
}

func (suite *AccountTestSuite) TestDeleteAccount() {
	var (
		account *storage.Account
		err     error
	)

	namespace := &storage.Namespace{
		Name: "myNamespace",
		Accounts: []*storage.Account{
			{Name: "Account1", Token: "token1"},
			{Name: "Account2", Token: "token2"},
			{Name: "Account3", Token: "token3"},
		},
	}

	suite.Len(namespace.Accounts, 3)
	account, err = namespace.FindAccount("Account1")
	suite.Require().NoError(err)

	namespace.DeleteAccount(account)
	suite.Len(namespace.Accounts, 2)
	account, err = namespace.FindAccount("Account1")
	suite.Require().ErrorIs(err, storage.NotFoundError{Type: "account", Name: "Account1"})
	// Delete again :D
	namespace.DeleteAccount(account)
	suite.Len(namespace.Accounts, 2)
}
