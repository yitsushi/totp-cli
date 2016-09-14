package main

import "errors"

type Namespace struct {
	Name     string
	Accounts []*Account
}

func (n *Namespace) FindAccount(name string) (account *Account, err error) {
	for _, account = range n.Accounts {
		if account.Name == name {
			return
		}
	}
	account = nil
	err = errors.New("Account not found.")

	return
}

func (n *Namespace) DeleteAccount(account *Account) {
	var position int = -1
	for i, item := range n.Accounts {
		if item == account {
			position = i
			break
		}
	}

	copy(n.Accounts[position:], n.Accounts[position+1:])
	n.Accounts[len(n.Accounts)-1] = nil
	n.Accounts = n.Accounts[:len(n.Accounts)-1]
}
