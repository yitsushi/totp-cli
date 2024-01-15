package main

import (
	"os"
)

func main() {
	app := newApplication()
	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}
