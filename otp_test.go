package main

import (
	"encoding/base32"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTOTP(t *testing.T) {
	secret := base32.StdEncoding.EncodeToString([]byte("82394783472398472348"))
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
		code := GenerateOTPCode(secret, when)
		assert.Equal(t, expected, code, when.String())
	}
}
