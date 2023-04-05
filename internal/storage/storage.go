package storage

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"

	"filippo.io/age"
	"filippo.io/age/armor"

	"github.com/yitsushi/totp-cli/internal/security"
	"github.com/yitsushi/totp-cli/internal/terminal"
)

const (
	dataLengthHead              = 13
	storageDirectoryPermissions = 0o700
	storageFilePermissions      = 0o600
)

// age defaults to a workFactor of 18 which works out to about a 1 second
// delay. The go docs (and the author of scrypt) suggest that a more
// interactive use case should target about 100ms. totp-cli uses a lower
// value of 15 in order to meet that goal. This can be adjusted in the
// future without breaking compatibility with existing files.
const scryptWorkFactor = 15

// Storage structure represents the credential storage.
type Storage struct {
	File     string `json:"-"`
	Password string `json:"-"`

	Namespaces []*Namespace
}

// DecryptV1 tries to decrypt the original unsecure SHA1 storage.
func (s *Storage) DecryptV1() error {
	encryptedData, err := os.ReadFile(s.File)
	if err != nil {
		return BackendError{Message: err.Error()}
	}

	decodedData, err := base64.StdEncoding.DecodeString(string(encryptedData))
	if err != nil {
		return BackendError{Message: err.Error()}
	}

	iv := decodedData[:aes.BlockSize]
	decodedData = decodedData[aes.BlockSize:]

	key := security.UnsecureSHA1(s.Password)

	block, err := aes.NewCipher(key)
	if err != nil {
		return BackendError{Message: err.Error()}
	}

	if len(decodedData)%aes.BlockSize != 0 {
		return BackendError{
			Message: "ciphertext is not a multiple of the block size",
		}
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(decodedData, decodedData)

	return s.parse(decodedData)
}

// DecryptV2 tries to decrypt the newer age storage.
func (s *Storage) DecryptV2() error {
	rawFile, err := os.Open(s.File)
	if err != nil {
		return BackendError{Message: err.Error()}
	}
	defer rawFile.Close()

	identity, err := age.NewScryptIdentity(s.Password)
	if err != nil {
		return BackendError{Message: err.Error()}
	}

	armorFile := armor.NewReader(rawFile)

	cryptFile, err := age.Decrypt(armorFile, identity)
	if err != nil {
		return BackendError{Message: err.Error()}
	}

	decodedData := &bytes.Buffer{}

	_, err = io.Copy(decodedData, cryptFile)
	if err != nil {
		return BackendError{Message: err.Error()}
	}

	return s.parse(decodedData.Bytes())
}

// Decrypt tries to decrypt the storage.
func (s *Storage) Decrypt() error {
	if err := s.DecryptV1(); err == nil {
		return nil
	}

	return s.DecryptV2()
}

// Save tries to encrypt and save the storage.
func (s *Storage) Save() error {
	tmpFile, err := os.CreateTemp(filepath.Dir(s.File),
		fmt.Sprintf("%s.*.tmp", filepath.Base(s.File)))
	if err != nil {
		return BackendError{Message: err.Error()}
	}

	defer os.Remove(tmpFile.Name())

	err = tmpFile.Chmod(storageFilePermissions)
	if err != nil {
		return BackendError{Message: err.Error()}
	}

	recipient, err := age.NewScryptRecipient(s.Password)
	if err != nil {
		return BackendError{Message: err.Error()}
	}

	recipient.SetWorkFactor(scryptWorkFactor)

	armorFile := armor.NewWriter(tmpFile)

	cryptFile, err := age.Encrypt(armorFile, recipient)
	if err != nil {
		return BackendError{Message: err.Error()}
	}

	err = json.NewEncoder(cryptFile).Encode(s.Namespaces)
	if err != nil {
		return BackendError{Message: err.Error()}
	}

	err = cryptFile.Close()
	if err != nil {
		return BackendError{Message: err.Error()}
	}

	err = armorFile.Close()
	if err != nil {
		return BackendError{Message: err.Error()}
	}

	err = tmpFile.Sync()
	if err != nil {
		return BackendError{Message: err.Error()}
	}

	err = tmpFile.Close()
	if err != nil {
		return BackendError{Message: err.Error()}
	}

	err = os.Rename(tmpFile.Name(), s.File)
	if err != nil {
		return BackendError{Message: err.Error()}
	}

	return nil
}

// FindNamespace returns with a namespace
// if the namespace does not exist error is not nil.
func (s *Storage) FindNamespace(name string) (*Namespace, error) {
	for _, namespace := range s.Namespaces {
		if namespace.Name == name {
			return namespace, nil
		}
	}

	return nil, NotFoundError{Type: "namespace", Name: name}
}

// DeleteNamespace removes a specific namespace from the Storage.
func (s *Storage) DeleteNamespace(namespace *Namespace) {
	position := -1

	for i, item := range s.Namespaces {
		if item == namespace {
			position = i

			break
		}
	}

	if position < 0 {
		return
	}

	copy(s.Namespaces[position:], s.Namespaces[position+1:])
	s.Namespaces[len(s.Namespaces)-1] = nil
	s.Namespaces = s.Namespaces[:len(s.Namespaces)-1]
}

// PrepareStorage loads, decrypt and parse the Storage. If the storage file does not exists it creates one.
func PrepareStorage() (*Storage, error) {
	credentialFile, storage, err := initStorage()
	if err != nil {
		return nil, err
	}

	if storage == nil {
		term := terminal.New(os.Stdin, os.Stdout, os.Stderr)

		password, termErr := term.Hidden("Password:")
		if termErr != nil {
			return nil, BackendError{Message: err.Error()}
		}

		storage = &Storage{
			File:     credentialFile,
			Password: password,
		}
	}

	err = storage.Decrypt()

	return storage, err
}

func initStorage() (string, *Storage, error) {
	var credentialFile string

	credentialFile = os.Getenv("TOTP_CLI_CREDENTIAL_FILE")

	if credentialFile == "" {
		currentUser, err := user.Current()
		if err != nil {
			return "", nil, BackendError{Message: err.Error()}
		}

		homePath := currentUser.HomeDir
		documentDirectory := filepath.Join(homePath, ".config/totp-cli")

		_, err = os.Stat(documentDirectory)
		if os.IsNotExist(err) {
			err = os.MkdirAll(documentDirectory, storageDirectoryPermissions)
		}

		if err != nil {
			return "", nil, BackendError{Message: err.Error()}
		}

		credentialFile = filepath.Join(documentDirectory, "credentials")
	}

	if _, err := os.Stat(credentialFile); err == nil {
		return credentialFile, nil, nil
	}

	term := terminal.New(os.Stdin, os.Stdout, os.Stderr)

	password, err := term.Hidden("Your Password (do not forget it):")
	if err != nil {
		return "", nil, BackendError{Message: err.Error()}
	}

	storage := &Storage{
		File:     credentialFile,
		Password: password,
	}

	err = storage.Save()

	return credentialFile, storage, err
}

func (s *Storage) parse(decodedData []byte) error {
	if err := s.parseV1(decodedData); err == nil {
		return nil
	}

	return s.parseV2(decodedData)
}

func (s *Storage) parseV1(decodedData []byte) error {
	var parsedData map[string]map[string]string

	// remove junk
	originalDataLength := bytes.IndexByte(decodedData, 0)
	if originalDataLength == 0 {
		originalDataLength = bytes.IndexByte(decodedData, dataLengthHead)
	}

	if originalDataLength > 0 && originalDataLength < len(decodedData) {
		decodedData = decodedData[:originalDataLength]
	}

	err := json.Unmarshal(decodedData, &parsedData)
	if err != nil {
		return BackendError{
			Message: "Something went wrong. Maybe this Password is not a valid one.",
		}
	}

	namespaces := []*Namespace{}

	for namespaceName, value := range parsedData {
		var accounts []*Account

		for accountName, secretKey := range value {
			account := &Account{Name: accountName, Token: secretKey}
			accounts = append(accounts, account)
		}

		namespace := &Namespace{Name: namespaceName, Accounts: accounts}
		namespaces = append(namespaces, namespace)
	}

	s.Namespaces = namespaces

	return nil
}

func (s *Storage) parseV2(decodedData []byte) error {
	namespaces := []*Namespace{}

	// remove junk
	originalDataLength := bytes.IndexByte(decodedData, 0)
	if originalDataLength == 0 {
		originalDataLength = bytes.IndexByte(decodedData, dataLengthHead)
	}

	if originalDataLength > 0 && originalDataLength < len(decodedData) {
		decodedData = decodedData[:originalDataLength]
	}

	err := json.Unmarshal(decodedData, &namespaces)
	if err != nil {
		return BackendError{
			Message: "Something went wrong. Maybe this Password is not a valid one.",
		}
	}

	s.Namespaces = namespaces

	return nil
}
