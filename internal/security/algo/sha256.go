package algo

import "crypto/sha256"

// SHA256 is a struct that implements the Algorithm interface for SHA256.
type SHA256 struct{}

// Hasher returns a function that returns a new hash.Hash.
func (s SHA256) Hasher() HasherFn {
	return sha256.New
}

var _ Algorithm = SHA256{}
