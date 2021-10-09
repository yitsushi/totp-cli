package security_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yitsushi/totp-cli/internal/security"
)

func TestUnsecureSHA1(t *testing.T) {
	table := map[string][]byte{
		"asdf":  []byte("3da541559918a808c2402bba5012f6c6"),
		"12345": []byte("8cb2237d0679ca88db6464eac60da963"),
	}

	for input, expected := range table {
		code := security.UnsecureSHA1(input)

		assert.Equal(t, expected, code)
	}
}
