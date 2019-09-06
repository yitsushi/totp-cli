package command

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"

	"github.com/kardianos/osext"
	grc "github.com/yitsushi/github-release-check"
	"github.com/yitsushi/go-commander"

	"github.com/yitsushi/totp-cli/info"
	"github.com/yitsushi/totp-cli/util"
)

// Update structure is the representation of the update command
type Update struct {
}

// Execute is the main function. It will be called on update command
func (c *Update) Execute(opts *commander.CommandHelper) {
	hasUpdate, release, _ := grc.Check(info.AppRepoOwner, info.AppName, info.AppVersion)

	if !hasUpdate {
		fmt.Printf("Your %s is up-to-date. \\o/\n", info.AppName)
		return
	}

	var assetToDownload *grc.Asset
	for _, asset := range release.Assets {
		if asset.Name == c.buildFilename(release.TagName) {
			assetToDownload = &asset
			break
		}
	}

	if assetToDownload == nil {
		fmt.Printf("Your %s is up-to-date. \\o/\n", info.AppName)
		return
	}

	c.downloadBinary(assetToDownload.BrowserDownloadURL)

	fmt.Printf("Now you have a fresh new %s \\o/\n", info.AppName)
}

func (c *Update) buildFilename(version string) string {
	return fmt.Sprintf("%s_%s_%s_%s.tar.gz", info.AppName, version, runtime.GOOS, runtime.GOARCH)
}

func (c *Update) downloadBinary(uri string) {
	fmt.Println("Download...")
	response, err := http.Get(uri)
	util.Check(err)
	defer response.Body.Close()

	gzipReader, _ := gzip.NewReader(response.Body)
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)
	tarReader.Next()

	file, _ := ioutil.TempFile("", info.AppName)
	util.Check(err)
	defer file.Close()

	_, err = io.Copy(file, tarReader)
	util.Check(err)

	file.Chmod(0755)

	currentExecutable, _ := osext.Executable()
	os.Rename(file.Name(), currentExecutable)
}

// NewUpdate creates a new Update command
func NewUpdate(appName string) *commander.CommandWrapper {
	return &commander.CommandWrapper{
		Handler: &Update{},
		Help: &commander.CommandDescriptor{
			Name:             "update",
			ShortDescription: fmt.Sprintf("Check and update %s itself", appName),
			LongDescription: `Check for updates.
If there is a newer version of this application for this OS and ARCH,
then download it and replace this application with the newer one.`,
		},
	}
}
