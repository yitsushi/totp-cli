package commander

// CommandInterface defined a command.
// If a struct implements all the required function,
// it is acceptable as a Command for CommandRegistry
type CommandInterface interface {
	Execute()
	Description() string
	ArgumentDescription() string
	Help() string
	Examples() []string
}
