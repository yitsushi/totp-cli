package main

import "flag"

const AppName string = "totp-cli"
const AppVersion string = "1.0.2"

func main() {
	flag.Parse()

	var command CommandFunction
	var ok bool

	if command, ok = commandHandlers[flag.Arg(0)]; !ok {
		command = Command_Help
	}

	command()
}
