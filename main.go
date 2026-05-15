package main

import (
	"fmt"
	"os"
)

func main() {
	app := newApplication()

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, " !! Error: %s\n", err.Error())

		os.Exit(1)
	}
}
