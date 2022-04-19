package main

import (
	"github.com/yitsushi/go-commander"

	"github.com/yitsushi/totp-cli/internal/cmd"
)

func registerCommands(registry *commander.CommandRegistry) {
	registry.Register(cmd.NewAddToken)
	registry.Register(cmd.NewChangePassword)
	registry.Register(cmd.NewDelete)
	registry.Register(cmd.NewDump)
	registry.Register(cmd.NewGenerate)
	registry.Register(cmd.NewImport)
	registry.Register(cmd.NewInstant)
	registry.Register(cmd.NewList)
	registry.Register(cmd.NewSetPrefix)
	registry.Register(cmd.NewUpdate)
	registry.Register(cmd.NewVersion)
}

func main() {
	registry := commander.NewCommandRegistry()

	registerCommands(registry)

	registry.Execute()
}
