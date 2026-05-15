package storage

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"io"
	"os"
	"os/user"
	"path/filepath"

	"filippo.io/age"
	"filippo.io/age/armor"
	"github.com/yitsushi/totp-cli/internal/security"
)

const (
	dataLengthHead              = 13
	storageDirectoryPermissions = 0o700
	storageFilePermissions      = 0o600
)

// age defaults to a workFactor of 18 which works out to about a second
// delay. The go docs (and the author of scrypt) suggest that a more
// interactive use case should target about 100ms. totp-cli uses a lower
// value of 15 in order to meet that goal. This can be adjusted in the
// future without breaking compatibility with existing files.
const scryptWorkFactor = 15

// FileBackend structure represents the credential storage.
type FileBackend struct {
	file             string `json:"-"`
	password         string `json:"-"`
	passwordProvider PasswordProvider

	namespaces []*Namespace
}

// NewFileStorage creates a new File Backend Storage.
func NewFileStorage(opts ...func(*FileBackend)) *FileBackend {
	backend := &FileBackend{}

	for _, o := range opts {
		o(backend)
	}

	return backend
}

// WithPasswordProvider sets the PasswordProvider used when a password is
// needed and TOTP_PASS is not set.
func WithPasswordProvider(pp PasswordProvider) func(*FileBackend) {
	return func(fb *FileBackend) {
		fb.passwordProvider = pp
	}
}

// SetPassword sets the password for the file storage.
func (s *FileBackend) SetPassword(password string) {
	s.password = password
}

// Prepare tries to load the credentials file and tries to decrypt it.
func (s *FileBackend) Prepare() error {
	err := s.initfileStorage()
	if err != nil {
		return err
	}

	if s.password != "" {
		return s.Decrypt()
	}

	s.password = os.Getenv("TOTP_PASS")

	if s.password == "" {
		s.password, err = s.acquirePassword("Password:")
		if err != nil {
			return err
		}
	}

	return s.Decrypt()
}

// Decrypt tries to decrypt the storage.
func (s *FileBackend) Decrypt() error {
	err := s.decryptV1()
	if err == nil {
		return nil
	}

	return s.decryptV2()
}

// FindNamespace returns with a namespace
// if the namespace does not exist error is not nil.
func (s *FileBackend) FindNamespace(name string) (*Namespace, error) {
	for _, namespace := range s.namespaces {
		if namespace.Name == name {
			return namespace, nil
		}
	}

	return nil, NotFoundError{Type: "namespace", Name: name}
}

// DeleteNamespace removes a specific namespace from the fileStorage.
func (s *FileBackend) DeleteNamespace(namespace *Namespace) {
	position := -1

	for i, item := range s.namespaces {
		if item == namespace {
			position = i

			break
		}
	}

	if position < 0 {
		return
	}

	copy(s.namespaces[position:], s.namespaces[position+1:])
	s.namespaces[len(s.namespaces)-1] = nil
	s.namespaces = s.namespaces[:len(s.namespaces)-1]
}

// AddNamespace adds a namespace to the namespace list if it's not already
// there.
func (s *FileBackend) AddNamespace(ns *Namespace) (*Namespace, error) {
	lookupNS, err := s.FindNamespace(ns.Name)
	if err == nil {
		return lookupNS, BackendError{
			Message: "namespace already exists: " + ns.Name,
		}
	}

	s.namespaces = append(s.namespaces, ns)

	return ns, nil
}

// ListNamespaces returns with all the namespaces defined in the storage.
func (s *FileBackend) ListNamespaces() []*Namespace {
	return s.namespaces
}

// Save tries to encrypt and save the storage.
func (s *FileBackend) Save() error {
	tmpFile, err := os.CreateTemp(
		filepath.Dir(s.file),
		filepath.Base(s.file)+".*.tmp",
	)
	if err != nil {
		return BackendError{Message: err.Error()}
	}

	defer os.Remove(tmpFile.Name())

	err = tmpFile.Chmod(storageFilePermissions)
	if err != nil {
		return BackendError{Message: err.Error()}
	}

	recipient, err := age.NewScryptRecipient(s.password)
	if err != nil {
		return BackendError{Message: err.Error()}
	}

	recipient.SetWorkFactor(scryptWorkFactor)

	armorFile := armor.NewWriter(tmpFile)

	cryptFile, err := age.Encrypt(armorFile, recipient)
	if err != nil {
		return BackendError{Message: err.Error()}
	}

	err = json.NewEncoder(cryptFile).Encode(s.namespaces)
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

	err = os.Rename(tmpFile.Name(), s.file)
	if err != nil {
		return BackendError{Message: err.Error()}
	}

	return nil
}

// acquirePassword asks the passwordProvider for a password using the given prompt.
func (s *FileBackend) acquirePassword(prompt string) (string, error) {
	if s.passwordProvider == nil {
		return "", BackendError{Message: "no password provider configured"}
	}

	password, err := s.passwordProvider.GetPassword(prompt)
	if err != nil {
		return "", BackendError{Message: err.Error()}
	}

	return password, nil
}

// decryptV1 tries to decrypt the storage with AES encryption using the SHA1
// hash of the password as encryption key.
func (s *FileBackend) decryptV1() error {
	encryptedData, err := os.ReadFile(s.file)
	if err != nil {
		return BackendError{Message: err.Error()}
	}

	decodedData, err := base64.StdEncoding.DecodeString(string(encryptedData))
	if err != nil {
		return BackendError{Message: err.Error()}
	}

	iv := decodedData[:aes.BlockSize]
	decodedData = decodedData[aes.BlockSize:]

	key := security.UnsecureSHA1(s.password)

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

// decryptV2 tries to decrypt the storage with "age" encryption using the
// password as it is as password.
func (s *FileBackend) decryptV2() error {
	rawFile, err := os.Open(s.file)
	if err != nil {
		return BackendError{Message: err.Error()}
	}
	defer rawFile.Close()

	identity, err := age.NewScryptIdentity(s.password)
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

func (s *FileBackend) initfileStorage() error {
	credentialFile, err := s.credentialsFilePath()
	if err != nil {
		return err
	}

	_, err = os.Stat(credentialFile)
	if err == nil {
		s.file = credentialFile

		return nil
	}

	password := os.Getenv("TOTP_PASS")

	if password == "" {
		password, err = s.acquirePassword("Your Password (do not forget it):")
		if err != nil {
			return err
		}
	}

	s.file = credentialFile
	s.password = password

	return s.Save()
}

func (s *FileBackend) parse(decodedData []byte) error {
	err := s.parseV1(decodedData)
	if err == nil {
		return nil
	}

	return s.parseV2(decodedData)
}

func (s *FileBackend) parseV1(decodedData []byte) error {
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

	namespaces := make([]*Namespace, 0, len(parsedData))

	for namespaceName, value := range parsedData {
		var accounts []*Account

		for accountName, secretKey := range value {
			account := &Account{Name: accountName, Token: secretKey}
			accounts = append(accounts, account)
		}

		namespace := &Namespace{Name: namespaceName, Accounts: accounts}
		namespaces = append(namespaces, namespace)
	}

	s.namespaces = namespaces

	return nil
}

func (s *FileBackend) parseV2(decodedData []byte) error {
	var namespaces []*Namespace

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

	s.namespaces = namespaces

	return nil
}

func (s *FileBackend) credentialsFilePath() (string, error) {
	filePath := os.Getenv("TOTP_CLI_CREDENTIAL_FILE")
	if filePath != "" {
		return filePath, nil
	}

	currentUser, err := user.Current()
	if err != nil {
		return "", BackendError{Message: err.Error()}
	}

	homePath := currentUser.HomeDir
	documentDirectory := filepath.Join(homePath, ".config", "totp-cli")

	_, err = os.Stat(documentDirectory)
	if os.IsNotExist(err) {
		err = os.MkdirAll(documentDirectory, storageDirectoryPermissions)
	}

	if err != nil {
		return "", BackendError{Message: err.Error()}
	}

	return filepath.Join(documentDirectory, "credentials"), nil
}
