package util

import (
	"fmt"
	"os"
)

// Check is for error handlig. If an error occurred it will simply
// Exit from the application with status code 1.
func Check(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
