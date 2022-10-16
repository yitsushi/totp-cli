package cmd

import "fmt"

// DownloadError is an error during downloading an update.
type DownloadError struct {
	Message string
}

func (e DownloadError) Error() string {
	return fmt.Sprintf("download error: %s", e.Message)
}

// ImportError is an error during a file import.
type ImportError struct {
	Message string
}

func (e ImportError) Error() string {
	return fmt.Sprintf("import error: %s", e.Message)
}

// GenerateError is an error during code generation.
type GenerateError struct {
	Message string
}

func (e GenerateError) Error() string {
	return fmt.Sprintf("generate error: %s", e.Message)
}

// DeleteError is an error during entry deletion.
type DeleteError struct {
	Message string
}

func (e DeleteError) Error() string {
	return fmt.Sprintf("delete error: %s", e.Message)
}
