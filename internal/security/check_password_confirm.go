package security

// CheckPasswordConfirm checks two byte array if the content is the same.
func CheckPasswordConfirm(password, confirm []byte) bool {
	if password == nil && confirm == nil {
		return true
	}

	if password == nil || confirm == nil {
		return false
	}

	if len(password) != len(confirm) {
		return false
	}

	for i := range password {
		if password[i] != confirm[i] {
			return false
		}
	}

	return true
}
