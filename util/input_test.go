package util

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type PasswordPair struct {
	Password []byte
	Confirm  []byte
	Correct  bool
}

func TestCheckPasswordConfirm(t *testing.T) {
	var passwordPairs []PasswordPair = []PasswordPair{
		PasswordPair{[]byte("asdf"), []byte("asdf"), true},
		PasswordPair{[]byte("asdfg"), []byte("asdf"), false},
		PasswordPair{[]byte("asdfg"), []byte("asdfh"), false},
		PasswordPair{[]byte("asdf"), []byte("asdfh"), false},
		PasswordPair{[]byte("asdf"), nil, false},
		PasswordPair{nil, []byte("asdf"), false},
		PasswordPair{nil, nil, true},
	}

	for _, pair := range passwordPairs {
		assert.Equal(
			t,
			CheckPasswordConfirm(pair.Password, pair.Confirm),
			pair.Correct,
			fmt.Sprintf("%s == %s -> %t", pair.Password, pair.Confirm, pair.Correct),
		)
	}
}
