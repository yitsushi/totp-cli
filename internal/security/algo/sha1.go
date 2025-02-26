package algo

import "crypto/sha1" //nolint:gosec // RFC-4226 defined SHA-1 as algorithm.

// SHA1 is a struct that implements the Algorithm interface for SHA1.
type SHA1 struct{}

// Hasher returns a function that returns a new hash.Hash.
func (s SHA1) Hasher() HasherFn {
	return sha1.New
}

var _ Algorithm = SHA1{}
