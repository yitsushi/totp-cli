package storage

// PasswordProvider supplies a password, optionally prompting the user.
type PasswordProvider interface {
	GetPassword(prompt string) (string, error)
}

// StaticPasswordProvider returns a fixed password without user interaction.
// Intended for use in tests.
type StaticPasswordProvider struct {
	Password string
}

// GetPassword returns the static password, ignoring the prompt.
func (p StaticPasswordProvider) GetPassword(_ string) (string, error) {
	return p.Password, nil
}
