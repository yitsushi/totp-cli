package security

import "fmt"

// OTPError is an error describing an error during OTP generation.
type OTPError struct {
	Message string
}

func (e OTPError) Error() string {
	return fmt.Sprintf("otp error: %s", e.Message)
}

// CryptoError is an error describing an error during cryptographic
// operations.
type CryptoError struct {
	Message string
}

func (e CryptoError) Error() string {
	return fmt.Sprintf("crypto error: %s", e.Message)
}
