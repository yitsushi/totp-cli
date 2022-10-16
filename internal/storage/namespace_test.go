package storage_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yitsushi/totp-cli/internal/storage"
)

func TestFindAccount(t *testing.T) {
	namespace := &storage.Namespace{
		Name: "mynamespace",
		Accounts: []*storage.Account{
			{Name: "Account1", Token: "token1"},
			{Name: "Account2", Token: "token2"},
			{Name: "Account3", Token: "token3"},
		},
	}

	account, err := namespace.FindAccount("Account1")

	assert.Equal(t, err, nil, "Error should be nil")
	assert.Equal(t, account.Name, "Account1", "Found account name should be Account1")
}

func TestFindAccount_NotFound(t *testing.T) {
	namespace := &storage.Namespace{
		Name: "mynamespace",
		Accounts: []*storage.Account{
			{Name: "Account1", Token: "token1"},
			{Name: "Account2", Token: "token2"},
			{Name: "Account3", Token: "token3"},
		},
	}

	account, err := namespace.FindAccount("AccountNotFound")

	assert.EqualError(
		t,
		err,
		"account not found: AccountNotFound",
		"Error should be 'account not found: AccountNotFound'",
	)
	assert.Nil(t, account)
}

func TestDeleteAccount(t *testing.T) {
	var (
		account *storage.Account
		err     error
	)

	namespace := &storage.Namespace{
		Name: "mynamespace",
		Accounts: []*storage.Account{
			{Name: "Account1", Token: "token1"},
			{Name: "Account2", Token: "token2"},
			{Name: "Account3", Token: "token3"},
		},
	}

	assert.Equal(t, len(namespace.Accounts), 3)
	account, err = namespace.FindAccount("Account1")
	assert.NoError(t, err)

	namespace.DeleteAccount(account)
	assert.Equal(t, len(namespace.Accounts), 2)
	account, err = namespace.FindAccount("Account1")
	assert.EqualError(
		t,
		err,
		"account not found: Account1",
		"Error should be 'account not found: Account1'",
	)
	// Delete again :D
	namespace.DeleteAccount(account)
	assert.Equal(t, len(namespace.Accounts), 2)
}
