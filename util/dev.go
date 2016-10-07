package util

import (
	"fmt"
	"os"
)

func Debugln(message string) {
	if os.Getenv("DEBUG") == "true" {
		os.Stderr.Write([]byte(fmt.Sprintf("[debug] %s\n", message)))
	}
}
