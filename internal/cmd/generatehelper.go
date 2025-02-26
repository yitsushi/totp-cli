package cmd

import (
	"fmt"
	"time"

	"github.com/yitsushi/totp-cli/internal/security"
	"github.com/yitsushi/totp-cli/internal/security/algo"
	s "github.com/yitsushi/totp-cli/internal/storage"
)

func formatCode(code string, remaining int64, showRemaining bool) string {
	if showRemaining {
		return fmt.Sprintf("%s (remaining time: %ds)", code, remaining)
	}

	return code
}

func generateCode(account *s.Account) (string, int64) {
	var algorithm algo.Algorithm

	switch account.Algorithm {
	case "sha1":
		algorithm = algo.SHA1{}
	case "sha256":
		algorithm = algo.SHA256{}
	case "sha512":
		algorithm = algo.SHA512{}
	default:
		algorithm = algo.Default{}
	}

	code, remaining, err := security.GenerateOTPCode(security.GenerateOptions{
		Token:      account.Token,
		When:       time.Now(),
		Length:     account.Length,
		Algorithm:  algorithm,
		TimePeriod: account.TimePeriod,
	})
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())

		return "", 0
	}

	if account.Prefix != "" {
		code = account.Prefix + code
	}

	return code, remaining
}
