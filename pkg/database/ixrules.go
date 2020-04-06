package database

import (
	"sort"
	"strings"

	"github.com/Suburbia-io/dashboard/pkg/errors"
	"github.com/Suburbia-io/dashboard/pkg/helpers/validators"
)

func sanitizeIX(includes, excludes []string) error {
	includeSet := map[string]bool{}
	for i := range includes {
		sanitized := validators.SanitizeToken(includes[i])
		if _, ok := includeSet[sanitized]; ok {
			return errors.IXRuleInvalidIncludes
		}

		if len(sanitized) == 0 {
			return errors.IXRuleInvalidIncludes
		}

		if strings.Contains(sanitized, ";") {
			return errors.IXRuleInvalidIncludes
		}

		includes[i] = sanitized
		includeSet[sanitized] = true
	}

	for i := range excludes {
		sanitized := validators.SanitizeToken(excludes[i])
		if _, ok := includeSet[sanitized]; ok {
			return errors.IXRuleInvalidExcludes
		}

		if len(sanitized) == 0 {
			return errors.IXRuleInvalidExcludes
		}

		if strings.Contains(sanitized, ";") {
			return errors.IXRuleInvalidExcludes
		}

		excludes[i] = sanitized
		includeSet[sanitized] = true
	}

	if len(includes) == 0 {
		return errors.IXRuleInvalidIncludes
	}

	sort.Strings(includes)
	sort.Strings(excludes)
	return nil
}
