package cmd

// DownloadError is an error during downloading an update.
type DownloadError struct {
	Message string
}

func (e DownloadError) Error() string {
	return "download error: " + e.Message
}

// CommandError is an error during downloading an update.
type CommandError struct {
	Message string
}

func (e CommandError) Error() string {
	return "error: " + e.Message
}

func resourceNotFoundError(name string) CommandError {
	return CommandError{Message: name + " does not exist"}
}

// FlagError is an error during flag parsing.
type FlagError struct {
	Message string
}

func (e FlagError) Error() string {
	return "flag error: " + e.Message
}

func invalidAlgorithmError(value string) FlagError {
	return FlagError{Message: "Invalid algorithm: " + value}
}
