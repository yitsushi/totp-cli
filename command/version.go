package command

import (
	"fmt"
	"runtime"

	"github.com/Yitsushi/totp-cli/info"
)

// Version structure is the representation of the Version command
type Version struct {
}

// Execute is the main function. It will be called on version command
func (c *Version) Execute() {
	fmt.Printf("%s %s (%s/%s)\n", info.AppName, info.AppVersion, runtime.GOOS, runtime.GOARCH)
}

// ArgumentDescription descripts the required and potential arguments
func (c *Version) ArgumentDescription() string {
	return ""
}

// Description will be displayed as Description (woooo) in the general help
func (c *Version) Description() string {
	return "Print current version of this application"
}

// Help is a general (human readable) command specific (long) help
func (c *Version) Help() string {
	return ""
}

// Examples lists a few example as array. Will be used in the command specific help
func (c *Version) Examples() []string {
	return []string{""}
}
