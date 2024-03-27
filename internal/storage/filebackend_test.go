package storage_test

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/yitsushi/totp-cli/internal/storage"
)

func TestFileBackend(t *testing.T) {
	suite.Run(t, &FileBackendTestSuite{})
}

type FileBackendTestSuite struct {
	suite.Suite
	storage storage.Storage
}

func (suite *FileBackendTestSuite) SetupTest() {
	suite.storage = storage.NewFileStorage()
}

func (suite *FileBackendTestSuite) TestFindNamespace() {
	suite.storage.AddNamespace(&storage.Namespace{Name: "Namespace1"})
	suite.storage.AddNamespace(&storage.Namespace{Name: "Namespace2"})
	suite.storage.AddNamespace(&storage.Namespace{Name: "Namespace3"})

	namespace, err := suite.storage.FindNamespace("Namespace1")

	suite.Require().NoError(err)
	suite.Equal("Namespace1", namespace.Name, "Found namespace name should be Namespace1")
}

func (suite *FileBackendTestSuite) TestNamespaceNotFound() {
	suite.storage.AddNamespace(&storage.Namespace{Name: "Namespace1"})
	suite.storage.AddNamespace(&storage.Namespace{Name: "Namespace2"})
	suite.storage.AddNamespace(&storage.Namespace{Name: "Namespace3"})

	namespace, err := suite.storage.FindNamespace("NamespaceNotFound")

	suite.Require().ErrorIs(err, storage.NotFoundError{Type: "namespace", Name: "NamespaceNotFound"})
	suite.Nil(namespace, "Namespace should be nil")
}

func (suite *FileBackendTestSuite) TestAddNamespace() {
	suite.storage.AddNamespace(&storage.Namespace{Name: "Namespace1"})
	suite.storage.AddNamespace(&storage.Namespace{Name: "Namespace2"})
	suite.storage.AddNamespace(&storage.Namespace{Name: "Namespace3"})

	namespace, err := suite.storage.AddNamespace(&storage.Namespace{Name: "Namespace4"})

	suite.Require().NoError(err)
	suite.Equal("Namespace4", namespace.Name)
}

func (suite *FileBackendTestSuite) TestAddAlreadyExistingNamespace() {
	suite.storage.AddNamespace(&storage.Namespace{Name: "Namespace1"})
	suite.storage.AddNamespace(&storage.Namespace{Name: "Namespace2"})
	suite.storage.AddNamespace(&storage.Namespace{Name: "Namespace3"})

	namespace, err := suite.storage.AddNamespace(&storage.Namespace{Name: "Namespace3"})

	suite.Require().ErrorIs(err, storage.BackendError{Message: "namespace already exists: Namespace3"})
	suite.Equal("Namespace3", namespace.Name)
}

func (suite *FileBackendTestSuite) TestDeleteNamespace() {
	var (
		namespace *storage.Namespace
		err       error
	)

	suite.storage.AddNamespace(&storage.Namespace{Name: "Namespace1"})
	suite.storage.AddNamespace(&storage.Namespace{Name: "Namespace2"})
	suite.storage.AddNamespace(&storage.Namespace{Name: "Namespace3"})

	suite.Len(suite.storage.ListNamespaces(), 3)
	namespace, err = suite.storage.FindNamespace("Namespace1")
	suite.Require().NoError(err)

	suite.storage.DeleteNamespace(namespace)
	suite.Len(suite.storage.ListNamespaces(), 2)
	namespace, err = suite.storage.FindNamespace("Namespace1")
	suite.Require().ErrorIs(err, storage.NotFoundError{Type: "namespace", Name: "Namespace1"})
	// Delete again :D
	suite.storage.DeleteNamespace(namespace)
	suite.Len(suite.storage.ListNamespaces(), 2)
}

func (suite *FileBackendTestSuite) TestDeleteNonExistingNamespace() {
	suite.storage.AddNamespace(&storage.Namespace{Name: "Namespace1"})
	suite.storage.AddNamespace(&storage.Namespace{Name: "Namespace2"})
	suite.storage.AddNamespace(&storage.Namespace{Name: "Namespace3"})

	suite.Len(suite.storage.ListNamespaces(), 3)
	suite.storage.DeleteNamespace(&storage.Namespace{Name: "Namespace4"})
	suite.Len(suite.storage.ListNamespaces(), 3)
}

func (suite *FileBackendTestSuite) TestReadWrite() {
	tmpDir, err := os.MkdirTemp("", "totp-cli-test-*")
	if err != nil {
		return
	}
	credsFilepath := path.Join(tmpDir, "credentials")

	defer func() {
		os.RemoveAll(tmpDir)
	}()

	os.Setenv("TOTP_PASS", "password")
	os.Setenv("TOTP_CLI_CREDENTIAL_FILE", credsFilepath)

	suite.storage.Prepare()
	suite.Empty(suite.storage.ListNamespaces())
	suite.storage.AddNamespace(&storage.Namespace{Name: "ns1"})
	suite.storage.AddNamespace(&storage.Namespace{Name: "ns2"})
	suite.Len(suite.storage.ListNamespaces(), 2)
	suite.storage.Save()

	newStorage := storage.NewFileStorage()
	newStorage.Prepare()
	suite.Len(newStorage.ListNamespaces(), 2)
}

func (suite *FileBackendTestSuite) TestInvalidPassword() {
	tmpDir, err := os.MkdirTemp("", "totp-cli-test-*")
	if err != nil {
		return
	}
	credsFilepath := path.Join(tmpDir, "credentials")

	defer func() {
		os.RemoveAll(tmpDir)
	}()

	os.Setenv("TOTP_PASS", "password")
	os.Setenv("TOTP_CLI_CREDENTIAL_FILE", credsFilepath)

	err = suite.storage.Prepare()
	suite.Require().NoError(err)
	suite.Empty(suite.storage.ListNamespaces())
	suite.storage.AddNamespace(&storage.Namespace{Name: "ns1"})
	suite.storage.AddNamespace(&storage.Namespace{Name: "ns2"})
	suite.Len(suite.storage.ListNamespaces(), 2)
	suite.storage.Save()

	newStorage := storage.NewFileStorage()
	newStorage.SetPassword("new password")
	err = newStorage.Prepare()
	suite.Require().ErrorIs(err, storage.BackendError{Message: "no identity matched any of the recipients"})
}
