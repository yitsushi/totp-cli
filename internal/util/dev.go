package util

import (
	"fmt"
	"os"
)

// Debugln is only for debugging. If DEBUG more is enabled (currently runtime)
// the given message will be printed to the Stderr
// with a NewLine character at the end of the line.
func Debugln(message string) {
	if os.Getenv("DEBUG") == "true" {
		os.Stderr.Write([]byte(fmt.Sprintf("[debug] %s\n", message)))
	}
}
