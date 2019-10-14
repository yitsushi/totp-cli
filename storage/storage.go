package storage

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	"github.com/yitsushi/totp-cli/util"
)

// Storage structure represents the credential storage
type Storage struct {
	File     string
	Password []byte

	Namespaces []*Namespace
}

// Decrypt tries to decrypt the storage
func (s *Storage) Decrypt() {
	encryptedData, err := ioutil.ReadFile(s.File)
	util.Check(err)
	decodedData, _ := base64.StdEncoding.DecodeString(string(encryptedData))

	iv := decodedData[:aes.BlockSize]
	decodedData = decodedData[aes.BlockSize:]

	block, err := aes.NewCipher(s.Password)
	util.Check(err)

	if len(decodedData)%aes.BlockSize != 0 {
		panic("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(decodedData, decodedData)

	s.parse(decodedData)
}

// Save tries to encrypt and save the storage
func (s *Storage) Save() {
	jsonStruct := map[string]map[string]string{}

	for _, namespace := range s.Namespaces {
		jsonStruct[namespace.Name] = map[string]string{}
		for _, account := range namespace.Accounts {
			jsonStruct[namespace.Name][account.Name] = account.Token
		}
	}

	plaintext, err := json.Marshal(jsonStruct)
	util.Check(err)

	missing := aes.BlockSize - (len(plaintext) % aes.BlockSize)
	padded := make([]byte, len(plaintext)+missing, len(plaintext)+missing)
	copy(padded[:], plaintext)
	plaintext = padded

	if len(plaintext)%aes.BlockSize != 0 {
		panic("plaintext is not a multiple of the block size")
	}

	block, err := aes.NewCipher(s.Password)
	util.Check(err)

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	encodedContent := base64.StdEncoding.EncodeToString(ciphertext)

	err = ioutil.WriteFile(s.File, []byte(encodedContent), 0644)
	util.Check(err)
}

// FindNamespace returns with a namespace
// if the namespace does not exist error is not nil
func (s *Storage) FindNamespace(name string) (namespace *Namespace, err error) {
	for _, namespace = range s.Namespaces {
		if namespace.Name == name {
			return
		}
	}
	namespace = &Namespace{}
	err = errors.New("Namespace not found")

	return
}

// DeleteNamespace removes a specific namespace from the Storage
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

var storage *Storage

// PrepareStorage loads, decrypt and parse the Storage. If the storage file does not exists it creates one.
func PrepareStorage() *Storage {
	credentialFile := initStorage()

	if storage != nil {
		return storage
	}

	password := util.AskPassword(32, "")

	storage = &Storage{
		File:     credentialFile,
		Password: password,
	}

	storage.Decrypt()

	return storage
}

func initStorage() string {
	var credentialFile string

	credentialFile = os.Getenv("TOTP_CLI_CREDENTIAL_FILE")

	if credentialFile == "" {
		currentUser, err := user.Current()
		util.Check(err)
		homePath := currentUser.HomeDir
		documentDirectory := filepath.Join(homePath, ".config/totp-cli")

		if _, err := os.Stat(documentDirectory); err != nil {
			if os.IsNotExist(err) {
				err = os.MkdirAll(documentDirectory, 0700)
				util.Check(err)
			} else {
				util.Check(err)
			}
		}

		credentialFile = filepath.Join(documentDirectory, "credentials")
	}

	if _, err := os.Stat(credentialFile); err == nil {
		return credentialFile
	}

	password := util.AskPassword(32, "Your Password (do not forget it)")
	storage = &Storage{
		File:     credentialFile,
		Password: password,
	}

	storage.Save()

	return credentialFile
}

func (s *Storage) parse(decodedData []byte) {
	var parsedData map[string]map[string]string

	// remove junk
	originalDataLength := bytes.IndexByte(decodedData, 0)
	if originalDataLength == 0 {
		originalDataLength = bytes.IndexByte(decodedData, 13)
	}

	if originalDataLength > 0 && originalDataLength < len(decodedData) {
		decodedData = decodedData[:originalDataLength]
	}

	err := json.Unmarshal(decodedData, &parsedData)
	if err != nil {
		fmt.Println("Something went wrong. Maybe this Password is not a valid one.")
		os.Exit(1)
	}

	var namespaces []*Namespace

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
}
