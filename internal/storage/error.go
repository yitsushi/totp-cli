package storage

import "fmt"

// NotFoundError is an error to tell the user the requested entity is not there.
type NotFoundError struct {
	Type string
	Name string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("%s not found: %s", e.Type, e.Name)
}

// StorageError is an error to tell the user something went wrong processing
// the backend storage.
type StorageError struct {
	Message string
}

func (e StorageError) Error() string {
	return fmt.Sprintf("storage error: %s", e.Message)
}
