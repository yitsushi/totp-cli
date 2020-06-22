package main

import (
	"github.com/yitsushi/go-commander"
	"github.com/yitsushi/totp-cli/command"
)

func registerCommands(registry *commander.CommandRegistry) {
	// Register available commands
	registry.Register(command.NewAddToken)
	registry.Register(command.NewChangePassword)
	registry.Register(command.NewDelete)
	registry.Register(command.NewGenerate)
	registry.Register(command.NewInstant)
	registry.Register(command.NewList)
	registry.Register(command.NewUpdate)
	registry.Register(command.NewVersion)
}

func main() {
	registry := commander.NewCommandRegistry()

	registerCommands(registry)

	registry.Execute()
}
