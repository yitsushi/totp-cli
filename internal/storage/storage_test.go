package storage_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	s "github.com/yitsushi/totp-cli/internal/storage"
)

func TestFindNamespace(t *testing.T) {
	storage := &s.Storage{
		Namespaces: []*s.Namespace{
			{Name: "Namespace1"},
			{Name: "Namespace2"},
			{Name: "Namespace3"},
		},
	}

	namespace, err := storage.FindNamespace("Namespace1")

	assert.Equal(t, err, nil, "Error should be nil")
	assert.Equal(t, namespace.Name, "Namespace1", "Found namespace name should be Namespace1")
}

func TestFindNamespace_NotFound(t *testing.T) {
	storage := &s.Storage{
		Namespaces: []*s.Namespace{
			{Name: "Namespace1"},
			{Name: "Namespace2"},
			{Name: "Namespace3"},
		},
	}

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

	storage := &s.Storage{
		Namespaces: []*s.Namespace{
			{Name: "Namespace1"},
			{Name: "Namespace2"},
			{Name: "Namespace3"},
		},
	}

	assert.Equal(t, len(storage.Namespaces), 3)
	namespace, err = storage.FindNamespace("Namespace1")
	assert.NoError(t, err)

	storage.DeleteNamespace(namespace)
	assert.Equal(t, len(storage.Namespaces), 2)
	namespace, err = storage.FindNamespace("Namespace1")
	assert.EqualError(
		t,
		err,
		"namespace not found: Namespace1",
		"Error should be 'namespace not found: Namespace1'")
	// Delete again :D
	storage.DeleteNamespace(namespace)
	assert.Equal(t, len(storage.Namespaces), 2)
}

func TestDecryptV1(t *testing.T) {
	decodedData := []byte{0xc6, 0x16, 0x1a, 0x9f, 0x03, 0x32, 0x18, 0x8d,
	                      0x9e, 0xde, 0x1c, 0x15, 0x50, 0xe3, 0x7e, 0x26,
	                      0x87, 0xa3, 0x28, 0xc4, 0x87, 0xc8, 0x73, 0x9c,
	                      0x66, 0xef, 0x94, 0x1d, 0x58, 0x4c, 0xed, 0xa5}
	storage := &s.Storage{
		Password: "joshua",
	}

	// try the wrong decrypter first
	err := storage.DecryptV2(decodedData)
	assert.NotNil(t, err)

	err = storage.DecryptV1(decodedData)
	assert.Nil(t, err)
	// the database is expected to be empty
	assert.Equal(t, len(storage.Namespaces), 0)
}

func TestDecryptV2(t *testing.T) {
	decodedData := []byte{0xd4, 0xbb, 0x47, 0x3a, 0x79, 0xba, 0xd2, 0x30,
	                      0x9d, 0x0f, 0x34, 0x91, 0x3b, 0xe7, 0xc5, 0xb4,
	                      0x1e, 0xe8, 0xeb, 0xab, 0x61, 0xca, 0x74, 0x0d,
	                      0x50, 0x7f, 0xe5, 0xce, 0x55, 0x3e, 0x44, 0xdb,
	                      0x4a, 0xea, 0x71, 0x3c, 0x2e, 0xa4, 0x36, 0x64,
	                      0xa8, 0x1d, 0x95, 0xb2, 0x14, 0x3f, 0x87, 0xbf}
	storage := &s.Storage{
		Password: "joshua",
	}

	// try the wrong decrypter first
	err := storage.DecryptV1(decodedData)
	assert.NotNil(t, err)

	err = storage.DecryptV2(decodedData)
	assert.Nil(t, err)
	// the database is expected to be empty
	assert.Equal(t, len(storage.Namespaces), 0)
}
