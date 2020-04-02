package validate

import (
	"regexp"

	"github.com/Suburbia-io/dashboard/pkg/errors"
)

var (
	slugRegexp = regexp.MustCompile(`^[a-z0-9_]+(?:-[a-z0-9_]+)*$`)
)

func Slug(slug string) error {
	if len(slug) < 2 || len(slug) > 40 {
		return errors.InvalidSlug
	}

	if !slugRegexp.MatchString(slug) {
		return errors.InvalidSlug
	}

	return nil
}
