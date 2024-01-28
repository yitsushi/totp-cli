package storage_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	s "github.com/yitsushi/totp-cli/internal/storage"
)

func TestFindNamespace(t *testing.T) {
	storage := s.NewFileStorage()

	storage.AddNamespace(&s.Namespace{Name: "Namespace1"})
	storage.AddNamespace(&s.Namespace{Name: "Namespace2"})
	storage.AddNamespace(&s.Namespace{Name: "Namespace3"})

	namespace, err := storage.FindNamespace("Namespace1")

	assert.Equal(t, err, nil, "Error should be nil")
	assert.Equal(t, namespace.Name, "Namespace1", "Found namespace name should be Namespace1")
}

func TestFindNamespace_NotFound(t *testing.T) {
	storage := s.NewFileStorage()

	storage.AddNamespace(&s.Namespace{Name: "Namespace1"})
	storage.AddNamespace(&s.Namespace{Name: "Namespace2"})
	storage.AddNamespace(&s.Namespace{Name: "Namespace3"})

	namespace, err := storage.FindNamespace("NamespaceNotFound")

	assert.EqualError(
		t,
		err,
		"namespace not found: NamespaceNotFound",
		"Error should be 'namespace not found: NamespaceNotFound'",
	)
	assert.Nil(t, namespace, "Namespace should be nil")
}

func TestDeleteNamespace(t *testing.T) {
	var (
		namespace *s.Namespace
		err       error
	)

	storage := s.NewFileStorage()

	storage.AddNamespace(&s.Namespace{Name: "Namespace1"})
	storage.AddNamespace(&s.Namespace{Name: "Namespace2"})
	storage.AddNamespace(&s.Namespace{Name: "Namespace3"})

	assert.Equal(t, len(storage.ListNamespaces()), 3)
	namespace, err = storage.FindNamespace("Namespace1")
	assert.NoError(t, err)

	storage.DeleteNamespace(namespace)
	assert.Equal(t, len(storage.ListNamespaces()), 2)
	namespace, err = storage.FindNamespace("Namespace1")
	assert.EqualError(
		t,
		err,
		"namespace not found: Namespace1",
		"Error should be 'namespace not found: Namespace1'")
	// Delete again :D
	storage.DeleteNamespace(namespace)
	assert.Equal(t, len(storage.ListNamespaces()), 2)
}
