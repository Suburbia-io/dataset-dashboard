package validators

import (
	"regexp"

	"github.com/Suburbia-io/dashboard/pkg/errors"
)

var TagTypePattern = regexp.MustCompile(`^[a-z]+[a-z0-9_]*$`)

func TagType(tagType string) (string, error) {
	if TagTypePattern.MatchString(tagType) {
		return tagType, nil
	}
	return tagType, errors.InvalidTagType
}
