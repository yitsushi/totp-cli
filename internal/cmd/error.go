package cmd

// DownloadError is an error during downloading an update.
type DownloadError struct {
	Message string
}

func (e DownloadError) Error() string {
	return "download error: %s" + e.Message
}

// CommandError is an error during downloading an update.
type CommandError struct {
	Message string
}

func (e CommandError) Error() string {
	return "error: %s" + e.Message
}

func resourceNotFoundError(name string) CommandError {
	return CommandError{Message: name + " does not exist"}
}
