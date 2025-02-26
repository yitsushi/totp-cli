package security_test

import (
	"encoding/base32"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/yitsushi/totp-cli/internal/security"
	"github.com/yitsushi/totp-cli/internal/security/algo"
	"github.com/yitsushi/totp-cli/internal/storage"
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
		code, _, err := security.GenerateOTPCode(security.GenerateOptions{
			Token: input,
			When:  when,
		})

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
		code, _, err := security.GenerateOTPCode(security.GenerateOptions{
			Token:     input,
			When:      when,
			Length:    8,
			Algorithm: algo.SHA1{},
		})

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
		code, _, err := security.GenerateOTPCode(security.GenerateOptions{
			Token:     input,
			When:      when,
			Algorithm: algo.SHA1{},
		})

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
		code, _, err := security.GenerateOTPCode(security.GenerateOptions{
			Token:     input,
			When:      when,
			Length:    storage.DefaultTokenLength,
			Algorithm: algo.SHA1{},
		})

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
		code, _, err := security.GenerateOTPCode(security.GenerateOptions{
			Token:     input,
			When:      when,
			Length:    storage.DefaultTokenLength,
			Algorithm: algo.SHA1{},
		})

		suite.Require().Error(err)
		suite.Equal(expected, code, when.String())
	}
}

func (suite *GenerateOTPCodeTestSuite) TestSHA256() {
	input := "JBSWY3DPEHPK3PXPJBSWY3DPEHPK3PXPJBSWY3DPEHPK3PXP"
	table := map[time.Time]string{
		time.Date(1970, 1, 1, 0, 0, 59, 0, time.UTC):     "598909",
		time.Date(2005, 3, 18, 1, 58, 29, 0, time.UTC):   "343094",
		time.Date(2005, 3, 18, 1, 58, 31, 0, time.UTC):   "342278",
		time.Date(2009, 2, 13, 23, 31, 30, 0, time.UTC):  "657794",
		time.Date(2016, 9, 16, 12, 40, 12, 0, time.UTC):  "139801",
		time.Date(2033, 5, 18, 3, 33, 20, 0, time.UTC):   "102968",
		time.Date(2603, 10, 11, 11, 33, 20, 0, time.UTC): "625152",
		time.Date(2025, 02, 26, 18, 12, 11, 0, time.UTC): "356698",
	}

	for when, expected := range table {
		code, _, err := security.GenerateOTPCode(security.GenerateOptions{
			Token:     input,
			When:      when,
			Length:    storage.DefaultTokenLength,
			Algorithm: algo.SHA256{},
		})

		suite.Require().NoError(err)
		suite.Equal(expected, code, when.String())
	}
}

func (suite *GenerateOTPCodeTestSuite) TestSHA512() {
	input := "JBSWY3DPEHPK3PXPJBSWY3DPEHPK3PXPJBSWY3DPEHPK3PXP"
	table := map[time.Time]string{
		time.Date(1970, 1, 1, 0, 0, 59, 0, time.UTC):     "735781",
		time.Date(2005, 3, 18, 1, 58, 29, 0, time.UTC):   "630426",
		time.Date(2005, 3, 18, 1, 58, 31, 0, time.UTC):   "719335",
		time.Date(2009, 2, 13, 23, 31, 30, 0, time.UTC):  "390343",
		time.Date(2016, 9, 16, 12, 40, 12, 0, time.UTC):  "760292",
		time.Date(2033, 5, 18, 3, 33, 20, 0, time.UTC):   "255524",
		time.Date(2603, 10, 11, 11, 33, 20, 0, time.UTC): "041274",
		time.Date(2025, 02, 26, 18, 12, 11, 0, time.UTC): "546487",
	}

	for when, expected := range table {
		code, _, err := security.GenerateOTPCode(security.GenerateOptions{
			Token:     input,
			When:      when,
			Length:    storage.DefaultTokenLength,
			Algorithm: algo.SHA512{},
		})

		suite.Require().NoError(err)
		suite.Equal(expected, code, when.String())
	}
}

func (suite *GenerateOTPCodeTestSuite) TestSHA256WithLongerPeriod() {
	input := "JBSWY3DPEHPK3PXPJBSWY3DPEHPK3PXPJBSWY3DPEHPK3PXP"
	table := map[time.Time]string{
		time.Date(2025, 02, 26, 18, 12, 1, 0, time.UTC):  "195134",
		time.Date(2025, 02, 26, 18, 12, 11, 0, time.UTC): "195134",
		time.Date(2025, 02, 26, 18, 12, 23, 0, time.UTC): "195134",
		time.Date(2025, 02, 26, 18, 12, 33, 0, time.UTC): "195134",
		time.Date(2025, 02, 26, 18, 12, 43, 0, time.UTC): "195134",
		time.Date(2025, 02, 26, 18, 12, 53, 0, time.UTC): "195134",
		time.Date(2025, 02, 26, 18, 13, 3, 0, time.UTC):  "042795",
		time.Date(2025, 02, 26, 18, 13, 13, 0, time.UTC): "042795",
	}

	for when, expected := range table {
		code, _, err := security.GenerateOTPCode(security.GenerateOptions{
			Token:      input,
			When:       when,
			Length:     storage.DefaultTokenLength,
			TimePeriod: 60,
			Algorithm:  algo.SHA256{},
		})

		suite.Require().NoError(err)
		suite.Equal(expected, code, when.String())
	}
}
