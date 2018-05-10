package main

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yitsushi/go-commander"
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
