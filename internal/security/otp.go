package security

import (
	"crypto/hmac"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/yitsushi/totp-cli/internal/security/algo"
)

const (
	mask1              = 0xf
	mask2              = 0x7f
	mask3              = 0xff
	passwordHashLength = 32
	shift16            = 16
	shift24            = 24
	shift8             = 8
	sumByteLength      = 8

	// DefaultLength is the default length of the generated OTP code.
	DefaultLength = 6
	// DefaultTimePeriod is the default time period for the TOTP.
	DefaultTimePeriod = 30
)

// GenerateOptions is the option list for the GenerateOTPCode function.
type GenerateOptions struct {
	Token      string
	When       time.Time
	Length     uint
	Algorithm  algo.Algorithm
	TimePeriod int64
}

func (opts *GenerateOptions) normalise() {
	if opts.Length == 0 {
		opts.Length = DefaultLength
	}

	if opts.Algorithm == nil {
		opts.Algorithm = algo.SHA1{}
	}

	if opts.TimePeriod == 0 {
		opts.TimePeriod = DefaultTimePeriod
	}

	// Remove spaces, some providers are giving us in a readable format,
	// so they add spaces in there. If it's not removed while pasting in,
	// remove it now.
	opts.Token = strings.ReplaceAll(opts.Token, " ", "")

	// It should be uppercase always
	opts.Token = strings.ToUpper(opts.Token)
}

// GenerateOTPCode generates an N digit TOTP from the secret Token.
func GenerateOTPCode(opts GenerateOptions) (string, int64, error) {
	opts.normalise()

	timer := uint64(math.Floor(float64(opts.When.Unix()) / float64(opts.TimePeriod)))
	remainingTime := opts.TimePeriod - opts.When.Unix()%opts.TimePeriod

	secretBytes, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(opts.Token)
	if err != nil {
		return "", 0, OTPError{Message: err.Error()}
	}

	buf := make([]byte, sumByteLength)
	mac := hmac.New(opts.Algorithm.Hasher(), secretBytes)

	binary.BigEndian.PutUint64(buf, timer)
	_, _ = mac.Write(buf)
	sum := mac.Sum(nil)

	// http://tools.ietf.org/html/rfc4226#section-5.4
	offset := sum[len(sum)-1] & mask1
	value := int64(((int(sum[offset]) & mask2) << shift24) |
		((int(sum[offset+1] & mask3)) << shift16) |
		((int(sum[offset+2] & mask3)) << shift8) |
		(int(sum[offset+3]) & mask3))

	//nolint:gosec // If the user sets a size that high to get an overflow, it's on them.
	modulo := int32(value % int64(math.Pow10(int(opts.Length))))

	format := fmt.Sprintf("%%0%dd", opts.Length)

	return fmt.Sprintf(format, modulo), remainingTime, nil
}
