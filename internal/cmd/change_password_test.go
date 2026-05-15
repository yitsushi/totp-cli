package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/yitsushi/totp-cli/internal/terminal"
)

func TestChangePassword(t *testing.T) {
	suite.Run(t, &ChangePasswordTestSuite{})
}

type ChangePasswordTestSuite struct {
	suite.Suite
}

func (suite *ChangePasswordTestSuite) TestChangePassword() {
	stor := newFakeStorage()

	err := executeChangePassword(stor, "newpassword")

	suite.Require().NoError(err)
	suite.True(stor.SaveCalled)
}

func (suite *ChangePasswordTestSuite) TestAskForNewPasswordMismatch() {
	term := terminal.New(
		bytes.NewReader([]byte("password1\npassword2\n")),
		&bytes.Buffer{},
		&bytes.Buffer{},
	)

	_, err := askForNewPassword(term)

	suite.Require().Error(err)
}

func (suite *ChangePasswordTestSuite) TestAskForNewPasswordMatch() {
	term := terminal.New(
		bytes.NewReader([]byte("mysecret\nmysecret\n")),
		&bytes.Buffer{},
		&bytes.Buffer{},
	)

	password, err := askForNewPassword(term)

	suite.Require().NoError(err)
	suite.Equal("mysecret", password)
}

func (suite *ChangePasswordTestSuite) TestAskForNewPasswordEOFReturnsError() {
	term := terminal.New(bytes.NewReader([]byte{}), &bytes.Buffer{}, &bytes.Buffer{})

	_, err := askForNewPassword(term)

	suite.Require().Error(err)
}

func (suite *ChangePasswordTestSuite) TestAskForNewPasswordEOFOnConfirmReturnsError() {
	term := terminal.New(bytes.NewReader([]byte("password\n")), &bytes.Buffer{}, &bytes.Buffer{})

	_, err := askForNewPassword(term)

	suite.Require().Error(err)
}
