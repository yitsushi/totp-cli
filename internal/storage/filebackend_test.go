package storage_test

import (
	"errors"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/yitsushi/totp-cli/internal/storage"
)

type errorPasswordProvider struct{}

func (e errorPasswordProvider) GetPassword(_ string) (string, error) {
	return "", errors.New("no tty")
}

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

	os.Setenv("TOTP_CLI_CREDENTIAL_FILE", credsFilepath)
	defer os.Unsetenv("TOTP_CLI_CREDENTIAL_FILE")

	stor := storage.NewFileStorage(
		storage.WithPasswordProvider(storage.StaticPasswordProvider{Password: "password"}),
	)

	err = stor.Prepare()
	suite.Require().NoError(err)
	suite.Empty(stor.ListNamespaces())
	stor.AddNamespace(&storage.Namespace{Name: "ns1"})
	stor.AddNamespace(&storage.Namespace{Name: "ns2"})
	suite.Len(stor.ListNamespaces(), 2)
	stor.Save()

	newStorage := storage.NewFileStorage(
		storage.WithPasswordProvider(storage.StaticPasswordProvider{Password: "password"}),
	)

	err = newStorage.Prepare()
	suite.Require().NoError(err)
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

	os.Setenv("TOTP_CLI_CREDENTIAL_FILE", credsFilepath)
	defer os.Unsetenv("TOTP_CLI_CREDENTIAL_FILE")

	stor := storage.NewFileStorage(
		storage.WithPasswordProvider(storage.StaticPasswordProvider{Password: "password"}),
	)

	err = stor.Prepare()
	suite.Require().NoError(err)
	suite.Empty(stor.ListNamespaces())
	stor.AddNamespace(&storage.Namespace{Name: "ns1"})
	stor.AddNamespace(&storage.Namespace{Name: "ns2"})
	suite.Len(stor.ListNamespaces(), 2)
	stor.Save()

	newStorage := storage.NewFileStorage(
		storage.WithPasswordProvider(storage.StaticPasswordProvider{Password: "wrong password"}),
	)

	err = newStorage.Prepare()
	suite.Require().ErrorIs(err, storage.BackendError{Message: "identity did not match any of the recipients: incorrect identity for recipient block: incorrect passphrase"})
}

func (suite *FileBackendTestSuite) TestSetPassword() {
	tmpDir, err := os.MkdirTemp("", "totp-cli-test-*")
	if err != nil {
		return
	}
	credsFilepath := path.Join(tmpDir, "credentials")

	defer func() {
		os.RemoveAll(tmpDir)
	}()

	os.Setenv("TOTP_CLI_CREDENTIAL_FILE", credsFilepath)
	defer os.Unsetenv("TOTP_CLI_CREDENTIAL_FILE")

	writer := storage.NewFileStorage(
		storage.WithPasswordProvider(storage.StaticPasswordProvider{Password: "mypassword"}),
	)
	suite.Require().NoError(writer.Prepare())

	reader := storage.NewFileStorage()
	reader.SetPassword("mypassword")

	err = reader.Prepare()
	suite.Require().NoError(err)
}

func (suite *FileBackendTestSuite) TestPrepareAcquirePasswordNilProvider() {
	tmpDir, err := os.MkdirTemp("", "totp-cli-test-*")
	if err != nil {
		return
	}
	credsFilepath := path.Join(tmpDir, "credentials")

	defer func() {
		os.RemoveAll(tmpDir)
	}()

	os.Setenv("TOTP_CLI_CREDENTIAL_FILE", credsFilepath)
	defer os.Unsetenv("TOTP_CLI_CREDENTIAL_FILE")
	os.Unsetenv("TOTP_PASS")

	writer := storage.NewFileStorage(
		storage.WithPasswordProvider(storage.StaticPasswordProvider{Password: "pass"}),
	)
	suite.Require().NoError(writer.Prepare())

	reader := storage.NewFileStorage()

	err = reader.Prepare()
	suite.Require().ErrorIs(err, storage.BackendError{Message: "no password provider configured"})
}

func (suite *FileBackendTestSuite) TestPrepareAcquirePasswordErrorProvider() {
	tmpDir, err := os.MkdirTemp("", "totp-cli-test-*")
	if err != nil {
		return
	}
	credsFilepath := path.Join(tmpDir, "credentials")

	defer func() {
		os.RemoveAll(tmpDir)
	}()

	os.Setenv("TOTP_CLI_CREDENTIAL_FILE", credsFilepath)
	defer os.Unsetenv("TOTP_CLI_CREDENTIAL_FILE")
	os.Unsetenv("TOTP_PASS")

	writer := storage.NewFileStorage(
		storage.WithPasswordProvider(storage.StaticPasswordProvider{Password: "pass"}),
	)
	suite.Require().NoError(writer.Prepare())

	reader := storage.NewFileStorage(
		storage.WithPasswordProvider(errorPasswordProvider{}),
	)

	err = reader.Prepare()
	suite.Require().ErrorIs(err, storage.BackendError{Message: "no tty"})
}
