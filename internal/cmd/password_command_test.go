package cmd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadPasswordCommand(t *testing.T) {
	password, err := readPasswordCommand("echo command-password")

	require.NoError(t, err)
	require.Equal(t, "command-password", password)
}

func TestReadPasswordCommandEmpty(t *testing.T) {
	_, err := readPasswordCommand(" ")

	require.EqualError(t, err, "password command is empty")
}
