package algo_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/yitsushi/totp-cli/internal/security/algo"
)

func TestAlgo(t *testing.T) {
	suite.Run(t, &AlgoTestSuite{})
}

type AlgoTestSuite struct {
	suite.Suite
}

func (suite *AlgoTestSuite) TestSHA1Hasher() {
	hasherFn := algo.SHA1{}.Hasher()
	suite.NotNil(hasherFn)
	suite.NotNil(hasherFn())
}

func (suite *AlgoTestSuite) TestSHA256Hasher() {
	hasherFn := algo.SHA256{}.Hasher()
	suite.NotNil(hasherFn)
	suite.NotNil(hasherFn())
}

func (suite *AlgoTestSuite) TestSHA512Hasher() {
	hasherFn := algo.SHA512{}.Hasher()
	suite.NotNil(hasherFn)
	suite.NotNil(hasherFn())
}

func (suite *AlgoTestSuite) TestDefaultHasher() {
	hasherFn := algo.Default{}.Hasher()
	suite.NotNil(hasherFn)
	suite.NotNil(hasherFn())
}
