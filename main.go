package main

import "flag"

func main() {
	flag.Parse()

	var command CommandFunction
	var ok bool

	if command, ok = commandHandlers[flag.Arg(0)]; !ok {
		command = Command_Help
	}

	command()
}
