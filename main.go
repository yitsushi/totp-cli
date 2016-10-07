package main

import (
	"github.com/Yitsushi/totp-cli/command"
	"github.com/Yitsushi/totp-cli/commander"
)

func registerCommands(registry *commander.CommandRegistry) {
	// Register available commands
	registry.Register("generate", &command.Generate{})
	registry.Register("add-token", &command.AddToken{})
	registry.Register("list", &command.List{})
	registry.Register("delete", &command.Delete{})
	registry.Register("change-password", &command.ChangePassword{})
	registry.Register("update", &command.Update{})
	registry.Register("version", &command.Version{})
}

func main() {
	registry := commander.NewCommandRegistry()

	registerCommands(registry)

	registry.Execute()
}
