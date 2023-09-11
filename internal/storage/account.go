package storage

// DefaultTokenLength defines what is the default length of the generated token.
// Most services are using 6 characters.
const DefaultTokenLength = 6

// Account represents a TOTP account.
type Account struct {
	Name   string
	Token  string
	Prefix string
	Length uint
}
