package security

// OTPError is an error describing an error during generation.
type OTPError struct {
	Message string
}

func (e OTPError) Error() string {
	return "otp error: " + e.Message
}
