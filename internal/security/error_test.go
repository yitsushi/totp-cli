package security_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/yitsushi/totp-cli/internal/security"
)

func TestSecurityError(t *testing.T) {
	suite.Run(t, &SecurityErrorTestSuite{})
}

type SecurityErrorTestSuite struct {
	suite.Suite
}

func (suite *SecurityErrorTestSuite) TestOTPErrorMessage() {
	err := security.OTPError{Message: "bad input"}
	suite.Equal("otp error: bad input", err.Error())
}
