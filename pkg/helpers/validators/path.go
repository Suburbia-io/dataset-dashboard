package validators

import (
	"regexp"

	"github.com/Suburbia-io/dashboard/pkg/errors"
)

var (
	pathRegexp = regexp.MustCompile(`^[a-z0-9\-]+(?:/[a-z0-9\-]+)*$`)
)

func Path(path string) (string, error) {
	if path == "" {
		return path, nil
	}

	if len(path) > 100 {
		return "", errors.InvalidPath
	}

	if !pathRegexp.MatchString(path) {
		return "", errors.InvalidPath
	}

	return path, nil
}
