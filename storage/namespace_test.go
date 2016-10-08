package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindAccount(t *testing.T) {
	namespace := &Namespace{
		Name: "mynamespace",
		Accounts: []*Account{
			&Account{Name: "Account1", Token: "token1"},
			&Account{Name: "Account2", Token: "token2"},
			&Account{Name: "Account3", Token: "token3"},
		},
	}

	account, err := namespace.FindAccount("Account1")

	assert.Equal(t, err, nil, "Error should be nil")
	assert.Equal(t, account.Name, "Account1", "Found account name should be Account1")
}

func TestFindAccount_NotFound(t *testing.T) {
	namespace := &Namespace{
		Name: "mynamespace",
		Accounts: []*Account{
			&Account{Name: "Account1", Token: "token1"},
			&Account{Name: "Account2", Token: "token2"},
			&Account{Name: "Account3", Token: "token3"},
		},
	}

	account, err := namespace.FindAccount("AccountNotFound")

	assert.EqualError(t, err, "Account not found.", "Error sould be 'Account not found.'")
	assert.Equal(t, account, &Account{}, "Account should be nil")
}

func TestDeleteAccount(t *testing.T) {
	var account *Account
	var err error

	namespace := &Namespace{
		Name: "mynamespace",
		Accounts: []*Account{
			&Account{Name: "Account1", Token: "token1"},
			&Account{Name: "Account2", Token: "token2"},
			&Account{Name: "Account3", Token: "token3"},
		},
	}

	assert.Equal(t, len(namespace.Accounts), 3)
	account, err = namespace.FindAccount("Account1")
	assert.Equal(t, err, nil, "Error should be nil")

	namespace.DeleteAccount(account)
	assert.Equal(t, len(namespace.Accounts), 2)
	account, err = namespace.FindAccount("Account1")
	assert.EqualError(t, err, "Account not found.", "Error sould be 'Account not found.'")
	// Delete again :D
	namespace.DeleteAccount(account)
	assert.Equal(t, len(namespace.Accounts), 2)
}
