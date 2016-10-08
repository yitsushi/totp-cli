package security

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/Yitsushi/totp-cli/util"
)

// GenerateOTPCode generates a 6 digit TOTP from the secret Token
func GenerateOTPCode(token string, when time.Time) string {
	timer := uint64(math.Floor(float64(when.Unix()) / float64(30)))
	secretBytes, err := base32.StdEncoding.DecodeString(strings.ToUpper(token))
	util.Check(err)

	buf := make([]byte, 8)
	mac := hmac.New(sha1.New, secretBytes)

	binary.BigEndian.PutUint64(buf, timer)
	mac.Write(buf)
	sum := mac.Sum(nil)

	// http://tools.ietf.org/html/rfc4226#section-5.4
	offset := sum[len(sum)-1] & 0xf
	value := int64(((int(sum[offset]) & 0x7f) << 24) |
		((int(sum[offset+1] & 0xff)) << 16) |
		((int(sum[offset+2] & 0xff)) << 8) |
		(int(sum[offset+3]) & 0xff))
	length := 6

	modulo := int32(value % int64(math.Pow10(length)))

	format := fmt.Sprintf("%%0%dd", length)
	return fmt.Sprintf(format, modulo)
}
