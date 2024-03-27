package security_test

import (
	"encoding/base32"
	"testing"
	"time"

	"github.com/yitsushi/totp-cli/internal/storage"

	"github.com/stretchr/testify/suite"

	"github.com/yitsushi/totp-cli/internal/security"
)

func TestGenerateOTPCodeSuit(t *testing.T) {
	suite.Run(t, &GenerateOTPCodeTestSuite{})
}

type GenerateOTPCodeTestSuite struct {
	suite.Suite
}

func (suite *GenerateOTPCodeTestSuite) TestDefault() {
	input := base32.StdEncoding.EncodeToString([]byte("82394783472398472348"))
	table := map[time.Time]string{
		time.Date(1970, 1, 1, 0, 0, 59, 0, time.UTC):     "007459",
		time.Date(2005, 3, 18, 1, 58, 29, 0, time.UTC):   "227921",
		time.Date(2005, 3, 18, 1, 58, 31, 0, time.UTC):   "638051",
		time.Date(2009, 2, 13, 23, 31, 30, 0, time.UTC):  "144100",
		time.Date(2016, 9, 16, 12, 40, 12, 0, time.UTC):  "346566",
		time.Date(2033, 5, 18, 3, 33, 20, 0, time.UTC):   "810915",
		time.Date(2603, 10, 11, 11, 33, 20, 0, time.UTC): "041334",
	}

	for when, expected := range table {
		code, _, err := security.GenerateOTPCode(input, when, storage.DefaultTokenLength)

		suite.Require().NoError(err)
		suite.Equal(expected, code, when.String())
	}
}

func (suite *GenerateOTPCodeTestSuite) TestDifferentLength() {
	input := base32.StdEncoding.EncodeToString([]byte("82394783472398472348"))
	table := map[time.Time]string{
		time.Date(1970, 1, 1, 0, 0, 59, 0, time.UTC):     "53007459",
		time.Date(2005, 3, 18, 1, 58, 29, 0, time.UTC):   "97227921",
		time.Date(2005, 3, 18, 1, 58, 31, 0, time.UTC):   "89638051",
		time.Date(2009, 2, 13, 23, 31, 30, 0, time.UTC):  "49144100",
		time.Date(2016, 9, 16, 12, 40, 12, 0, time.UTC):  "13346566",
		time.Date(2033, 5, 18, 3, 33, 20, 0, time.UTC):   "44810915",
		time.Date(2603, 10, 11, 11, 33, 20, 0, time.UTC): "28041334",
	}

	for when, expected := range table {
		code, _, err := security.GenerateOTPCode(input, when, 8)

		suite.Require().NoError(err)
		suite.Equal(expected, code, when.String())
	}
}

func (suite *GenerateOTPCodeTestSuite) TestSpaceSeparatedToken() {
	input := "37kh vdxt c5hj ttfp ujok cipy jy"
	table := map[time.Time]string{
		time.Date(1970, 1, 1, 0, 0, 59, 0, time.UTC):     "066634",
		time.Date(2005, 3, 18, 1, 58, 29, 0, time.UTC):   "597310",
		time.Date(2005, 3, 18, 1, 58, 31, 0, time.UTC):   "174182",
		time.Date(2009, 2, 13, 23, 31, 30, 0, time.UTC):  "623746",
		time.Date(2016, 9, 16, 12, 40, 12, 0, time.UTC):  "330739",
		time.Date(2033, 5, 18, 3, 33, 20, 0, time.UTC):   "556617",
		time.Date(2603, 10, 11, 11, 33, 20, 0, time.UTC): "608345",
	}

	for when, expected := range table {
		code, _, err := security.GenerateOTPCode(input, when, storage.DefaultTokenLength)

		suite.Require().NoError(err)
		suite.Equal(expected, code, when.String())
	}
}

func (suite *GenerateOTPCodeTestSuite) TestNonPaddedHashes() {
	input := "a6mryljlbufszudtjdt42nh5by"
	table := map[time.Time]string{
		time.Date(1970, 1, 1, 0, 0, 59, 0, time.UTC):     "866149",
		time.Date(2005, 3, 18, 1, 58, 29, 0, time.UTC):   "996077",
		time.Date(2005, 3, 18, 1, 58, 31, 0, time.UTC):   "421761",
		time.Date(2009, 2, 13, 23, 31, 30, 0, time.UTC):  "903464",
		time.Date(2016, 9, 16, 12, 40, 12, 0, time.UTC):  "997249",
		time.Date(2033, 5, 18, 3, 33, 20, 0, time.UTC):   "210476",
		time.Date(2603, 10, 11, 11, 33, 20, 0, time.UTC): "189144",
	}

	for when, expected := range table {
		code, _, err := security.GenerateOTPCode(input, when, storage.DefaultTokenLength)

		suite.Require().NoError(err)
		suite.Equal(expected, code, when.String())
	}
}

func (suite *GenerateOTPCodeTestSuite) TestInvalidPadding() {
	input := "a6mr*&^&*%*&ylj|'[lbufszudtjdt42nh5by"
	table := map[time.Time]string{
		time.Date(1970, 1, 1, 0, 0, 59, 0, time.UTC):   "",
		time.Date(2005, 3, 18, 1, 58, 29, 0, time.UTC): "",
	}

	for when, expected := range table {
		code, _, err := security.GenerateOTPCode(input, when, storage.DefaultTokenLength)

		suite.Require().Error(err)
		suite.Equal(expected, code, when.String())
	}
}
