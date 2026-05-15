package cmd

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/yitsushi/totp-cli/internal/storage"
)

func TestGenerate(t *testing.T) {
	suite.Run(t, &GenerateTestSuite{})
}

type GenerateTestSuite struct {
	suite.Suite
}

func (suite *GenerateTestSuite) TestGetAccountFound() {
	ns := &storage.Namespace{
		Name:     "myns",
		Accounts: []*storage.Account{{Name: "github", Token: "JBSWY3DPEHPK3PXP"}},
	}
	stor := newFakeStorage(ns)

	acc, err := getAccount(stor, "myns", "github")

	suite.Require().NoError(err)
	suite.Equal("github", acc.Name)
}

func (suite *GenerateTestSuite) TestGetAccountNamespaceNotFound() {
	stor := newFakeStorage()

	_, err := getAccount(stor, "missing", "github")

	suite.Require().Error(err)
}

func (suite *GenerateTestSuite) TestGetAccountAccountNotFound() {
	ns := &storage.Namespace{Name: "myns"}
	stor := newFakeStorage(ns)

	_, err := getAccount(stor, "myns", "missing")

	suite.Require().Error(err)
}

func (suite *GenerateTestSuite) TestFormatCodeWithoutRemaining() {
	result := formatCode("123456", 20, false)

	suite.Equal("123456", result)
}

func (suite *GenerateTestSuite) TestFormatCodeWithRemaining() {
	result := formatCode("123456", 20, true)

	suite.Equal("123456 (remaining time: 20s)", result)
}

func (suite *GenerateTestSuite) TestGenerateCodeSHA1() {
	acc := &storage.Account{Name: "github", Token: "JBSWY3DPEHPK3PXP", Algorithm: "sha1"}

	code, _, err := generateCode(acc)

	suite.Require().NoError(err)
	suite.Len(code, 6)
}

func (suite *GenerateTestSuite) TestGenerateCodeSHA256() {
	acc := &storage.Account{Name: "github", Token: "JBSWY3DPEHPK3PXP", Algorithm: "sha256"}

	code, _, err := generateCode(acc)

	suite.Require().NoError(err)
	suite.Len(code, 6)
}

func (suite *GenerateTestSuite) TestGenerateCodeSHA512() {
	acc := &storage.Account{Name: "github", Token: "JBSWY3DPEHPK3PXP", Algorithm: "sha512"}

	code, _, err := generateCode(acc)

	suite.Require().NoError(err)
	suite.Len(code, 6)
}

func (suite *GenerateTestSuite) TestGenerateCodeDefaultAlgorithm() {
	acc := &storage.Account{Name: "github", Token: "JBSWY3DPEHPK3PXP", Algorithm: ""}

	code, _, err := generateCode(acc)

	suite.Require().NoError(err)
	suite.Len(code, 6)
}

func (suite *GenerateTestSuite) TestGenerateCodeWithPrefix() {
	acc := &storage.Account{
		Name:   "github",
		Token:  "JBSWY3DPEHPK3PXP",
		Prefix: "GH:",
	}

	code, _, err := generateCode(acc)

	suite.Require().NoError(err)
	suite.True(strings.HasPrefix(code, "GH:"))
}

func (suite *GenerateTestSuite) TestGenerateCodeInvalidToken() {
	acc := &storage.Account{
		Name:  "github",
		Token: "!!!invalid!!!",
	}

	_, _, err := generateCode(acc)

	suite.Require().Error(err)
}

func (suite *GenerateTestSuite) TestExecuteGenerateSuccess() {
	ns := &storage.Namespace{
		Name:     "myns",
		Accounts: []*storage.Account{{Name: "github", Token: "JBSWY3DPEHPK3PXP"}},
	}
	stor := newFakeStorage(ns)

	err := executeGenerate(stor, "myns", "github", false, false)

	suite.Require().NoError(err)
}

func (suite *GenerateTestSuite) TestExecuteGenerateNamespaceNotFound() {
	stor := newFakeStorage()

	err := executeGenerate(stor, "missing", "github", false, false)

	suite.Require().Error(err)
}
