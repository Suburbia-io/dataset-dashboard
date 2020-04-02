package validators

import (
	"regexp"
	"strings"

	"github.com/Suburbia-io/dashboard/pkg/errors"
)

var (
	brandsBadSpacesRegexp   = regexp.MustCompile(`\r\n|\r|\n|\t`)
	brandsDoubleSpaceRegexp = regexp.MustCompile(`\s\s+`)
)

func BrandLabel(label string) (string, error) {
	label = strings.TrimSpace(label)

	if label == "" {
		return "", errors.BrandInvalidLabel
	}

	label = brandsDoubleSpaceRegexp.ReplaceAllString(label, " ")

	for _, r := range []*regexp.Regexp{brandsBadSpacesRegexp} {
		if r.MatchString(label) {
			return "", errors.BrandInvalidLabel
		}
	}

	return label, nil
}

func BrandDescription(description string) (string, error) {
	description = strings.TrimSpace(description)

	description = brandsDoubleSpaceRegexp.ReplaceAllString(description, " ")

	for _, r := range []*regexp.Regexp{brandsBadSpacesRegexp} {
		if r.MatchString(description) {
			return "", errors.BrandInvalidDescription
		}
	}

	return description, nil
}
