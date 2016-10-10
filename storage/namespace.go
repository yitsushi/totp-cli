package storage

import "errors"

// Namespace represents a Namespace "category"
type Namespace struct {
	Name     string
	Accounts []*Account
}

// FindAccount returns with an account under a specific Namespace
// if the account does not exist error is not nil
func (n *Namespace) FindAccount(name string) (account *Account, err error) {
	for _, account = range n.Accounts {
		if account.Name == name {
			return
		}
	}
	account = &Account{}
	err = errors.New("Account not found.")

	return
}

// DeleteAccount removes a specific Account from the Namespace
func (n *Namespace) DeleteAccount(account *Account) {
	position := -1
	for i, item := range n.Accounts {
		if item == account {
			position = i
			break
		}
	}

	if position < 0 {
		return
	}

	copy(n.Accounts[position:], n.Accounts[position+1:])
	n.Accounts[len(n.Accounts)-1] = nil
	n.Accounts = n.Accounts[:len(n.Accounts)-1]
}
