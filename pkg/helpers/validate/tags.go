package validate

import (
	"regexp"

	"github.com/Suburbia-io/dashboard/pkg/errors"
)

var TagTypePattern = regexp.MustCompile(`^[a-z]+[a-z0-9_]*$`)

func TagType(tagType string) error {
	if TagTypePattern.MatchString(tagType) {
		return nil
	}
	return errors.InvalidTagType
}

// Same regex as the path regex with the addition of +.
var tagRegexp = regexp.MustCompile(`^[a-z0-9\-+]+(?:/[a-z0-9\-+]+)*$`)

func Tag(tag string) error {
	if tag == "" {
		return errors.InvalidTag
	}

	if len(tag) > 100 {
		return errors.InvalidTag
	}

	if !tagRegexp.MatchString(tag) {
		return errors.InvalidTag
	}

	return nil
}

func TagGrade(g int) error {
	if g < 0 || g > 10 {
		return errors.InvalidConfidence
	}
	return nil
}
