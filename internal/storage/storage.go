package storage

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"os"
	"os/user"
	"path/filepath"

	"github.com/yitsushi/totp-cli/internal/security"
	"github.com/yitsushi/totp-cli/internal/terminal"
)

const (
	saltSize                    = 16
	dataLengthHead              = 13
	storageDirectoryPermissions = 0o700
	storageFilePermissions      = 0o600
)

// Storage structure represents the credential storage.
type Storage struct {
	File     string `json:"-"`
	Password string `json:"-"`

	Namespaces []*Namespace
}

// DecryptV1 tries to decrypt the original unsecure SHA1 storage.
//
// The AES key is from UnsecureSHA1(password)
//
// On disk format is a base64 encoded copy of:
// +--------+--------+------------------------+
// | Offset | Length | Value                  |
// +--------+--------+------------------------+
// | 0      | 16     | Initialization Vector  |
// +--------+--------+------------------------+
// | 16     | EOF    | AES-256 encrypted data |
// +--------+--------+------------------------+.
func (s *Storage) DecryptV1(decodedData []byte) error {
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
	plaintext := make([]byte, len(decodedData))
	mode.CryptBlocks(plaintext, decodedData)

	return s.parse(plaintext)
}

// DecryptV2 tries to decrypt the stronger scrypt storage.
//
// The AES key is from Scrypt(password, salt, 1<<17, 8, 1)
//
// On disk format is a base64 encoded copy of:
// +--------+--------+------------------------+
// | Offset | Length | Value                  |
// +--------+--------+------------------------+
// | 0      | 16     | Salt                   |
// +--------+--------+------------------------+
// | 16     | 16     | Initialization Vector  |
// +--------+--------+------------------------+
// | 32     | EOF    | AES-256 encrypted data |
// +--------+--------+------------------------+.
func (s *Storage) DecryptV2(decodedData []byte) error {
	salt := decodedData[:saltSize]
	decodedData = decodedData[saltSize:]
	iv := decodedData[:aes.BlockSize]
	decodedData = decodedData[aes.BlockSize:]

	key, err := security.Scrypt(s.Password, salt)
	if err != nil {
		return BackendError{Message: err.Error()}
	}

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
	plaintext := make([]byte, len(decodedData))
	mode.CryptBlocks(plaintext, decodedData)

	return s.parse(plaintext)
}

// Decrypt tries to decrypt the storage.
func (s *Storage) Decrypt() error {
	encryptedData, err := os.ReadFile(s.File)
	if err != nil {
		return BackendError{Message: err.Error()}
	}

	decodedData, err := base64.StdEncoding.DecodeString(string(encryptedData))
	if err != nil {
		return BackendError{Message: err.Error()}
	}

	if err := s.DecryptV2(decodedData); err == nil {
		return nil
	}

	return s.DecryptV1(decodedData)
}

// Save tries to encrypt and save the storage.
func (s *Storage) Save() error {
	plaintext, err := json.Marshal(s.Namespaces)
	if err != nil {
		return BackendError{Message: err.Error()}
	}

	missing := aes.BlockSize - (len(plaintext) % aes.BlockSize)
	padded := make([]byte, len(plaintext)+missing)

	copy(padded, plaintext)

	plaintext = padded

	if len(plaintext)%aes.BlockSize != 0 {
		return BackendError{Message: "plaintext is not a multiple of the block size"}
	}

	ciphertext := make([]byte, saltSize+aes.BlockSize+len(plaintext))

	salt := ciphertext[:saltSize]
	if _, readErr := io.ReadFull(rand.Reader, salt); readErr != nil {
		return BackendError{Message: readErr.Error()}
	}

	key, err := security.Scrypt(s.Password, salt)
	if err != nil {
		return BackendError{Message: err.Error()}
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return BackendError{Message: err.Error()}
	}

	iv := ciphertext[saltSize : saltSize+aes.BlockSize]
	if _, readErr := io.ReadFull(rand.Reader, iv); readErr != nil {
		return BackendError{Message: readErr.Error()}
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[saltSize+aes.BlockSize:], plaintext)

	encodedContent := base64.StdEncoding.EncodeToString(ciphertext)

	err = os.WriteFile(s.File, []byte(encodedContent), storageFilePermissions)
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
