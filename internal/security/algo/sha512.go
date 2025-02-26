package algo

import "crypto/sha512"

// SHA512 is a struct that implements the Algorithm interface for SHA512.
type SHA512 struct{}

// Hasher returns a function that returns a new hash.Hash.
func (s SHA512) Hasher() HasherFn {
	return sha512.New
}

var _ Algorithm = SHA512{}
