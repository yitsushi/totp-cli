package algo

import "hash"

// HasherFn is a function that returns a new hash.Hash.
type HasherFn func() hash.Hash

// Algorithm is an interface that defines the methods that an algorithm must implement.
type Algorithm interface {
	// Hasher returns a function that returns a new hash.Hash.
	Hasher() HasherFn
}
