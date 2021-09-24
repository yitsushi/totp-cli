package security

import (
	"crypto/hmac"
	"os"

	//nolint:gosec // yolo?
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"math"
	"strings"
	"time"
)

const (
	mask1              = 0xf
	mask2              = 0x7f
	mask3              = 0xff
	timeSplitInSeconds = 30
	shift24            = 24
	shift16            = 16
	shift8             = 8
	codeLength         = 6
	sumByteLength      = 8
)

// GenerateOTPCode generates a 6 digit TOTP from the secret Token.
func GenerateOTPCode(token string, when time.Time) string {
	timer := uint64(math.Floor(float64(when.Unix()) / float64(timeSplitInSeconds)))
	// Remove spaces, some providers are giving us in a readable format
	// so they add spaces in there. If it's not removed while pasting in,
	// remove it now.
	token = strings.ReplaceAll(token, " ", "")

	// It should be uppercase always
	token = strings.ToUpper(token)

	secretBytes, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(token)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
		os.Exit(1)
	}

	buf := make([]byte, sumByteLength)
	mac := hmac.New(sha1.New, secretBytes)

	binary.BigEndian.PutUint64(buf, timer)
	_, _ = mac.Write(buf)
	sum := mac.Sum(nil)

	// http://tools.ietf.org/html/rfc4226#section-5.4
	offset := sum[len(sum)-1] & mask1
	value := int64(((int(sum[offset]) & mask2) << shift24) |
		((int(sum[offset+1] & mask3)) << shift16) |
		((int(sum[offset+2] & mask3)) << shift8) |
		(int(sum[offset+3]) & mask3))

	modulo := int32(value % int64(math.Pow10(codeLength)))

	format := fmt.Sprintf("%%0%dd", codeLength)

	return fmt.Sprintf(format, modulo)
}
