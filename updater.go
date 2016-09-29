package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"

	grc "github.com/Yitsushi/github-release-check"
	"github.com/kardianos/osext"
)

type Updater struct {
}

func (u *Updater) CheckAndDownloadVersion() bool {
	hasUpdate, release, _ := grc.Check("yitsushi", AppName, AppVersion)

	if !hasUpdate {
		return false
	}

	var assetToDownload *grc.Asset
	for _, asset := range release.Assets {
		if asset.Name == u.BuildFilename() {
			assetToDownload = &asset
			break
		}
	}

	if assetToDownload == nil {
		return false
	}

	u.DownloadBinary(assetToDownload.BrowserDownloadURL)

	return true
}

func (u *Updater) BuildFilename() string {
	extension := ""
	if runtime.GOOS == "windows" {
		extension = ".exe"
	}
	return fmt.Sprintf("%s-%s-%s%s", AppName, runtime.GOOS, runtime.GOARCH, extension)
}

func (u *Updater) DownloadBinary(uri string) {
	fmt.Println("Download...")
	response, err := http.Get(uri)
	check(err)
	defer response.Body.Close()

	file, _ := ioutil.TempFile("", AppName)
	check(err)
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	check(err)

	file.Chmod(0755)

	currentExecutable, _ := osext.Executable()
	os.Rename(file.Name(), currentExecutable)
}
