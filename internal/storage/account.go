package storage

// Account represents a TOTP account.
type Account struct {
	Name   string
	Token  string
	Prefix string `yaml:"Prefix,omitempty"`
}
