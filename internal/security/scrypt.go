package security

import (
	"golang.org/x/crypto/scrypt"
)

// scrypt parameters.
const (
	costFactor            = 17
	blockSizeFactor       = 8
	parallelizationFactor = 1
)

// Wrapper around scrypt.Key() that ensures the use of a consistent set of
// hardening parameters.
func Scrypt(text string, salt []byte) ([]byte, error) {
	data, err := scrypt.Key([]byte(text), salt, 1<<costFactor, blockSizeFactor,
		parallelizationFactor, passwordHashLength)
	if err != nil {
		return nil, CryptoError{Message: err.Error()}
	}

	return data, nil
}
