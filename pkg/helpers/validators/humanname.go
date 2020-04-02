package validators

import (
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/Suburbia-io/dashboard/pkg/errors"
)

var (
	badSpacesRegexp    = regexp.MustCompile(`\r\n|\r|\n|\t`)
	doubleSpacesRegexp = regexp.MustCompile(`\s\s+`)
)

// HumanName sanitizes and validate a human's name.
//
// Sanitize:
//   * Remove leading and trailing spaces.
//   * Replace double spaces with single spaces.
//
// Validate:
//   * Check length (2-200)
//   * Check for tabs, newlines, etc.
func HumanName(name string) (string, error) {
	name = strings.TrimSpace(name)

	if strlen := utf8.RuneCountInString(name); strlen < 1 || strlen > 200 {
		return "", errors.InvalidHumanName
	}

	name = doubleSpacesRegexp.ReplaceAllString(name, " ")

	for _, r := range []*regexp.Regexp{badSpacesRegexp} {
		if r.MatchString(name) {
			return "", errors.InvalidHumanName
		}
	}

	return name, nil
}
