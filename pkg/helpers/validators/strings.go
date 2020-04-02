package validators

import (
	"strings"

	"github.com/Suburbia-io/dashboard/pkg/errors"
)

func NonEmptyString(s string) (string, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return s, errors.EmptyString
	}
	return s, nil
}

func TrimmedString(s string) (string, error) {
	return strings.TrimSpace(s), nil
}

func TrimmedNullableString(s *string) (*string, error) {
	if s == nil {
		return nil, nil
	}
	t := strings.TrimSpace(*s)
	return &t, nil
}
