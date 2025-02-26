package main

import (
	"fmt"
	"os"
)

func main() {
	app := newApplication()
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, " !! Error: %s\n", err.Error())

		os.Exit(1)
	}
}
