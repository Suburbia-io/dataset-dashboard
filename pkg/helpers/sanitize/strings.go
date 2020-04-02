package sanitize

import (
	"regexp"
	"strings"
)

var (
	doubleSpacesRegexp = regexp.MustCompile(`\s\s+`)
	badSpacesRegexp    = regexp.MustCompile(`\r\n|\r|\n|\t`)
)

func TrimSpace(s string) string {
	return strings.TrimSpace(s)
}

func SingleLineString(s string) string {
	s = badSpacesRegexp.ReplaceAllString(s, " ")
	s = doubleSpacesRegexp.ReplaceAllString(s, " ")
	s = strings.TrimSpace(s)
	return s
}

func SingleLineStringPtr(s *string) *string {
	if s == nil {
		return nil
	}
	t := SingleLineString(*s)
	return &t
}
