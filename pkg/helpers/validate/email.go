package validate

import (
	"net/mail"
	"strings"

	"github.com/Suburbia-io/dashboard/pkg/errors"
)

func Email(email string) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return errors.InvalidEmail
	}

	// Don't allow email addresses without a full domain - e.g. bill@localhost is not allowed.
	if domain := strings.Split(email, "@")[1]; !strings.Contains(domain, ".") {
		return errors.InvalidEmail
	}

	return nil
}
