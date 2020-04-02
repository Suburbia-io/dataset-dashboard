package validate

import (
	"github.com/Suburbia-io/dashboard/pkg/errors"
	zxcvbn "github.com/nbutton23/zxcvbn-go"
)

func StrongPassword(password string, userInputs []string) (err error) {
	result := zxcvbn.PasswordStrength(password, userInputs)
	if result.Score < 4 {
		return errors.PasswordTooWeak
	}
	return nil
}

func StrongPasswordScore(password string, userInputs []string) int {
	return zxcvbn.PasswordStrength(password, userInputs).Score
}
