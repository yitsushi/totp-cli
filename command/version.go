package command

import (
	"fmt"
	"runtime"

	"github.com/Yitsushi/totp-cli/info"
)

type Version struct {
}

func (c *Version) Execute() {
	fmt.Printf("%s %s (%s/%s)\n", info.AppName, info.AppVersion, runtime.GOOS, runtime.GOARCH)
}

func (c *Version) ArgumentDescription() string {
	return ""
}

func (c *Version) Description() string {
	return "Print current version of this application"
}

func (c *Version) Help() string {
	return ""
}

func (c *Version) Examples() []string {
	return []string{""}
}
