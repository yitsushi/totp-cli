package commander

type CommandInterface interface {
	Execute()
	Description() string
	ArgumentDescription() string
	Help() string
	Examples() []string
}
