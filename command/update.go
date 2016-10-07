package command

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"

	grc "github.com/Yitsushi/github-release-check"
	"github.com/kardianos/osext"

	"github.com/Yitsushi/totp-cli/info"
	"github.com/Yitsushi/totp-cli/util"
)

type Update struct {
}

func (c *Update) Description() string {
	return fmt.Sprintf("Check and update %s itself", info.AppName)
}

func (c *Update) ArgumentDescription() string {
	return ""
}

func (c *Update) Execute() {
	hasUpdate, release, _ := grc.Check(info.AppRepoOwner, info.AppName, info.AppVersion)

	if !hasUpdate {
		fmt.Printf("Your %s is up-to-date. \\o/\n", info.AppName)
		return
	}

	var assetToDownload *grc.Asset
	for _, asset := range release.Assets {
		if asset.Name == c.BuildFilename() {
			assetToDownload = &asset
			break
		}
	}

	if assetToDownload == nil {
		fmt.Printf("Your %s is up-to-date. \\o/\n", info.AppName)
		return
	}

	c.DownloadBinary(assetToDownload.BrowserDownloadURL)

	fmt.Printf("Now you have a fresh new %s \\o/\n", info.AppName)
}

func (c *Update) BuildFilename() string {
	extension := ""
	if runtime.GOOS == "windows" {
		extension = ".exe"
	}
	return fmt.Sprintf("%s-%s-%s%s", info.AppName, runtime.GOOS, runtime.GOARCH, extension)
}

func (c *Update) DownloadBinary(uri string) {
	fmt.Println("Download...")
	response, err := http.Get(uri)
	util.Check(err)
	defer response.Body.Close()

	file, _ := ioutil.TempFile("", info.AppName)
	util.Check(err)
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	util.Check(err)

	file.Chmod(0755)

	currentExecutable, _ := osext.Executable()
	os.Rename(file.Name(), currentExecutable)
}

func (c *Update) Help() string {
	return `Check for updates.
If there is a newer version of this application for this OS and ARCH,
then download it and replace this application with the newer one.`
}

func (c *Update) Examples() []string {
	return []string{""}
}
