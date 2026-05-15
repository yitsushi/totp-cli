package security

import "crypto/subtle"

// CheckPasswordConfirm checks two byte array if the content is the same.
func CheckPasswordConfirm(password, confirm []byte) bool {
	if password == nil && confirm == nil {
		return true
	}

	if password == nil || confirm == nil {
		return false
	}

	return subtle.ConstantTimeCompare(password, confirm) == 1
}
