package main

import (
	"flag"
	"os"
	"testing"

	"github.com/Yitsushi/totp-cli/commander"
	"github.com/stretchr/testify/assert"
)

func TestRegisterCommands(t *testing.T) {
	registry := commander.NewCommandRegistry()

	registerCommands(registry)

	assert.NotEqual(t, 0, len(registry.Commands))
}

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}
