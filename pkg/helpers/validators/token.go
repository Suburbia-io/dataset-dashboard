package validators

import (
	"regexp"
	"strings"
)

var spacesPtn = regexp.MustCompile(`\s+`)

func SanitizeToken(token string) string {
	return strings.ToLower(strings.TrimSpace(spacesPtn.ReplaceAllString(token, " ")))
}
