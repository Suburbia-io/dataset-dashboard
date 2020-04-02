package validators

import (
	"net/mail"

	"github.com/Suburbia-io/dashboard/pkg/errors"
)

// Email sanitizes and validates an email address.
//
// Sanitize:
//   * Make lowercase.
//   * Remove non-address information.
//
// Validate:
//   * Attempt to parse using Go's mail package.
func Email(email string) (string, error) {
	address, err := mail.ParseAddress(email)
	if err != nil {
		return "", errors.InvalidEmail
	}

	return address.Address, nil
}
