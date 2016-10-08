package commander

import (
	"flag"
	"fmt"

	"github.com/Yitsushi/totp-cli/info"
	"github.com/Yitsushi/totp-cli/util"
)

// CommandRegistry will handle all CLI request
// and find the route to the proper Command
type CommandRegistry struct {
	Commands map[string]CommandInterface

	maximumCommandLength int
}

// Register function is used to register a Command in CommandRegistry
// the first argument will be the command from the CLI
// the second argument will be the handler that implements CommandInterface
func (c *CommandRegistry) Register(name string, handler CommandInterface) {
	c.Commands[name] = handler
	commandLength := len(fmt.Sprintf("%s %s", name, handler.ArgumentDescription()))
	if commandLength > c.maximumCommandLength {
		c.maximumCommandLength = commandLength
	}
	util.Debugln(fmt.Sprintf("'%s' command is registered.", name))
}

// Execute finds the proper command, handle errors from the command and print Help
// if the given command it unknown or print the Command specific help
// if something went wrong or the user asked for it.
func (c *CommandRegistry) Execute() {
	name := flag.Arg(0)
	if command, ok := c.Commands[name]; ok {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("[E] %s\n\n", err)
				c.CommandHelp(name)
			}
		}()

		command.Execute()
	} else {
		c.Help()
	}
}

// Help lists all available commands to the user
func (c *CommandRegistry) Help() {
	if flag.Arg(0) == "help" && flag.Arg(1) != "" {
		c.CommandHelp(flag.Arg(1))
		return
	}

	format := fmt.Sprintf("%%-%ds   %%s\n", c.maximumCommandLength)
	for name, command := range c.Commands {
		fmt.Printf(
			format,
			fmt.Sprintf("%s %s", name, command.ArgumentDescription()),
			command.Description(),
		)
	}
}

// CommandHelp prints more detailed help for a specific Command
func (c *CommandRegistry) CommandHelp(name string) {
	util.Debugln(name)
	if command, ok := c.Commands[name]; ok {
		fmt.Printf("Usage: %s %s %s\n", info.AppName, name, command.ArgumentDescription())

		if command.Help() != "" {
			fmt.Println("")
			fmt.Println(command.Help())
		}

		if len(command.Examples()) > 0 {
			fmt.Printf("\nExamples:\n")
			for _, line := range command.Examples() {
				fmt.Printf("  %s %s %s\n", info.AppName, name, line)
			}
		}
	}
}

// NewCommandRegistry is a simple "constructor"-like function
// that initializes Commands map
func NewCommandRegistry() *CommandRegistry {
	flag.Parse()
	return &CommandRegistry{
		Commands: map[string]CommandInterface{},
	}
}
