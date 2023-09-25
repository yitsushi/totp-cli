package cmd

import (
	"fmt"
	"time"

	"github.com/yitsushi/totp-cli/internal/security"
	s "github.com/yitsushi/totp-cli/internal/storage"
)

func formatCode(code string, remaining int64, showRemaining bool) string {
	if showRemaining {
		return fmt.Sprintf("%s (remaining time: %ds)", code, remaining)
	}

	return code
}

func generateCode(account *s.Account) (string, int64) {
	code, remaining, err := security.GenerateOTPCode(account.Token, time.Now(), account.Length)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())

		return "", 0
	}

	if account.Prefix != "" {
		code = account.Prefix + code
	}

	return code, remaining
}
