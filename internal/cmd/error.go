package cmd

import "fmt"

// DownloadError is an error during downloading an update.
type DownloadError struct {
	Message string
}

func (e DownloadError) Error() string {
	return fmt.Sprintf("download error: %s", e.Message)
}

// CommandError is an error during downloading an update.
type CommandError struct {
	Message string
}

func (e CommandError) Error() string {
	return fmt.Sprintf("error: %s", e.Message)
}
