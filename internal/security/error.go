package security

import "fmt"

// OTPError is an error describing an error during generation.
type OTPError struct {
	Message string
}

func (e OTPError) Error() string {
	return fmt.Sprintf("otp error: %s", e.Message)
}
