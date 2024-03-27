package storage

// Namespace represents a Namespace "category".
type Namespace struct {
	Name     string     `json:"name"     yaml:"name"`
	Accounts []*Account `json:"accounts" yaml:"accounts"`
}

// FindAccount returns with an account under a specific Namespace
// if the account does not exist error is not nil.
func (n *Namespace) FindAccount(name string) (*Account, error) {
	for _, account := range n.Accounts {
		if account.Name == name {
			return account, nil
		}
	}

	return nil, NotFoundError{Type: "account", Name: name}
}

// DeleteAccount removes a specific Account from the Namespace.
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
