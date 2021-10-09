package cmd_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yitsushi/totp-cli/internal/cmd"
)

type PasswordPair struct {
	Password []byte
	Confirm  []byte
	Correct  bool
}

func TestCheckPasswordConfirm(t *testing.T) {
	passwordPairs := []PasswordPair{
		{[]byte("asdf"), []byte("asdf"), true},
		{[]byte("asdfg"), []byte("asdf"), false},
		{[]byte("asdfg"), []byte("asdfh"), false},
		{[]byte("asdf"), []byte("asdfh"), false},
		{[]byte("asdf"), nil, false},
		{nil, []byte("asdf"), false},
		{nil, nil, true},
	}

	for _, pair := range passwordPairs {
		assert.Equal(
			t,
			cmd.CheckPasswordConfirm(pair.Password, pair.Confirm),
			pair.Correct,
			fmt.Sprintf("%s == %s -> %t", pair.Password, pair.Confirm, pair.Correct),
		)
	}
}
