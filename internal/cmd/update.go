package cmd

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"runtime"

	"github.com/kardianos/osext"
	"github.com/urfave/cli/v2"
	grc "github.com/yitsushi/github-release-check"

	"github.com/yitsushi/totp-cli/internal/info"
)

func UpdateCommand() *cli.Command {
	return &cli.Command{
		Name:      "update",
		Usage:     fmt.Sprintf("Check and update %s itself", info.AppName),
		ArgsUsage: " ",
		Description: `Check for updates.
If there is a newer version of this application for this OS and ARCH,
then download it and replace this application with the newer one.`,
		Action: func(_ *cli.Context) error {
			hasUpdate, release, _ := grc.Check(info.AppRepoOwner, info.AppName, info.AppVersion)

			if !hasUpdate {
				fmt.Printf("Your %s is up-to-date. \\o/\n", info.AppName)

				return nil
			}

			var (
				assetToDownload grc.Asset
				found           bool
			)

			for _, asset := range release.Assets {
				buildFilename := fmt.Sprintf("%s-%s-%s-%s.tar.gz", info.AppName, release.TagName, runtime.GOOS, runtime.GOARCH)
				if asset.Name == buildFilename {
					assetToDownload = asset
					found = true

					break
				}
			}

			if !found {
				fmt.Printf("Your %s is up-to-date. \\o/\n", info.AppName)

				return nil
			}

			fmt.Printf("Target: %s\n", assetToDownload.Name)

			err := downloadBinary(assetToDownload.BrowserDownloadURL)
			if err != nil {
				return err
			}

			fmt.Printf("Now you have a fresh new %s \\o/\n", info.AppName)

			return nil
		},
	}
}

func downloadBinary(uri string) error {
	fmt.Println(" -> Download...")

	client := http.Client{}

	request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, uri, nil)
	if err != nil {
		return DownloadError{Message: err.Error()}
	}

	response, err := client.Do(request)
	if err != nil {
		return DownloadError{Message: err.Error()}
	}

	defer response.Body.Close()

	gzipReader, _ := gzip.NewReader(response.Body)
	defer gzipReader.Close()

	fmt.Println(" -> Extract...")

	tarReader := tar.NewReader(gzipReader)

	_, err = tarReader.Next()
	if err != nil {
		return DownloadError{Message: err.Error()}
	}

	currentExecutable, _ := osext.Executable()
	originalPath := path.Dir(currentExecutable)

	file, err := os.CreateTemp(originalPath, info.AppName)
	if err != nil {
		return DownloadError{Message: err.Error()}
	}

	defer file.Close()

	_, err = io.Copy(file, tarReader) //nolint:gosec // I don't have better option right now.
	if err != nil {
		return DownloadError{Message: err.Error()}
	}

	err = file.Chmod(binaryChmodValue)
	if err != nil {
		return DownloadError{Message: err.Error()}
	}

	err = os.Rename(file.Name(), currentExecutable)
	if err != nil {
		return DownloadError{Message: err.Error()}
	}

	return nil
}
