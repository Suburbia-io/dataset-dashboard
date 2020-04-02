package validators

import (
	"regexp"

	"github.com/Suburbia-io/dashboard/pkg/errors"
)

var sessionTokenPattern = regexp.MustCompile(`^[a-zA-Z0-9]{32}$`)

func SessionToken(token string) error {
	if !sessionTokenPattern.MatchString(string(token)) {
		return errors.InvalidSessionToken
	}
	return nil
}
