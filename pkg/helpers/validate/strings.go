package validate

import (
	"github.com/Suburbia-io/dashboard/pkg/errors"
)

func NonEmptyString(s string) error {
	if s == "" {
		return errors.EmptyString
	}
	return nil
}
