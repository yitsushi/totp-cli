package cmd

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestCmdError(t *testing.T) {
	suite.Run(t, &CmdErrorTestSuite{})
}

type CmdErrorTestSuite struct {
	suite.Suite
}

func (suite *CmdErrorTestSuite) TestDownloadErrorMessage() {
	err := DownloadError{Message: "timeout"}
	suite.Equal("download error: timeout", err.Error())
}

func (suite *CmdErrorTestSuite) TestCommandErrorMessage() {
	err := CommandError{Message: "bad"}
	suite.Equal("error: bad", err.Error())
}

func (suite *CmdErrorTestSuite) TestFlagErrorMessage() {
	err := FlagError{Message: "invalid"}
	suite.Equal("flag error: invalid", err.Error())
}

func (suite *CmdErrorTestSuite) TestInvalidAlgorithmError() {
	err := invalidAlgorithmError("md5")
	suite.Equal("flag error: Invalid algorithm: md5", err.Error())
}

func (suite *CmdErrorTestSuite) TestResourceNotFoundError() {
	err := resourceNotFoundError("myns")
	suite.Equal("error: myns does not exist", err.Error())
}
