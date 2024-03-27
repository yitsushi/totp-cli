package storage

// DefaultTokenLength defines what is the default length of the generated token.
// Most services are using 6 characters.
const DefaultTokenLength = 6

// Account represents a TOTP account.
type Account struct {
	Name   string `json:"name"   yaml:"name"`
	Token  string `json:"token"  yaml:"token"`
	Prefix string `json:"prefix" yaml:"prefix"`
	Length uint   `json:"length" yaml:"length"`
}
