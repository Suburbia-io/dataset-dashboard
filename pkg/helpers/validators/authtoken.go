package validators

import (
	"regexp"

	"github.com/Suburbia-io/dashboard/pkg/errors"
	"github.com/Suburbia-io/dashboard/pkg/helpers/crypto"
)

var authTokenPattern = regexp.MustCompile("^" +
	"[" + string(crypto.ReadableChars) + "]{4}-" +
	"[" + string(crypto.ReadableChars) + "]{4}-" +
	"[" + string(crypto.ReadableChars) + "]{4}-" +
	"[" + string(crypto.ReadableChars) + "]{4}$")

func AuthToken(token string) error {
	if !authTokenPattern.MatchString(string(token)) {
		return errors.InvalidAuthToken
	}
	return nil
}
