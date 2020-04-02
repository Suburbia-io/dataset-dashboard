package validate

import (
	"regexp"

	"github.com/Suburbia-io/dashboard/pkg/errors"
)

var (
	pathRegexp = regexp.MustCompile(`^[a-z0-9\-]+(?:/[a-z0-9\-]+)*$`)
)

func Path(path string) error {
	if path == "" {
		return nil
	}

	if len(path) > 100 {
		return errors.InvalidPath
	}

	if !pathRegexp.MatchString(path) {
		return errors.InvalidPath
	}

	return nil
}
