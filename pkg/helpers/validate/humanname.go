package validate

import (
	"unicode/utf8"

	"github.com/Suburbia-io/dashboard/pkg/errors"
)

// HumanName sanitizes and validate a human's name.
// Validate:
//   * Check length (2-200)
func HumanName(name string) error {
	if strlen := utf8.RuneCountInString(name); strlen < 1 || strlen > 200 {
		return errors.InvalidHumanName
	}
	return nil
}
