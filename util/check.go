package util

import (
	"fmt"
	"os"
)

func Check(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
