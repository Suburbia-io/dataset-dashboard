package validators

import (
	"regexp"
	"strings"

	"github.com/Suburbia-io/dashboard/pkg/errors"
)

var (
	corpBadSpacesRegexp = regexp.MustCompile(`\r\n|\r|\n|\t`)
	doubleSpaceRegexp   = regexp.MustCompile(`\s\s+`)
)

func CorpName(name string) (string, error) {
	name = strings.TrimSpace(name)

	name = doubleSpaceRegexp.ReplaceAllString(name, " ")

	for _, r := range []*regexp.Regexp{corpBadSpacesRegexp} {
		if r.MatchString(name) {
			return "", errors.CorporationInvalidName
		}
	}

	return name, nil
}

func CorpExchange(exchange string) (string, error) {
	exchange = strings.TrimSpace(exchange)

	exchange = doubleSpaceRegexp.ReplaceAllString(exchange, " ")

	for _, r := range []*regexp.Regexp{corpBadSpacesRegexp} {
		if r.MatchString(exchange) {
			return "", errors.CorporationInvalidSymbol
		}
	}

	return exchange, nil
}

func CorpCode(code string) (string, error) {
	code = strings.TrimSpace(code)

	code = doubleSpaceRegexp.ReplaceAllString(code, " ")

	for _, r := range []*regexp.Regexp{corpBadSpacesRegexp} {
		if r.MatchString(code) {
			return "", errors.CorporationInvalidSymbol
		}
	}

	return code, nil
}

func CorpISIN(isin string) (string, error) {
	isin = strings.TrimSpace(isin)

	isin = doubleSpaceRegexp.ReplaceAllString(isin, " ")

	for _, r := range []*regexp.Regexp{corpBadSpacesRegexp} {
		if r.MatchString(isin) {
			return "", errors.CorporationInvalidISIN
		}
	}

	return isin, nil
}

func CorpCUSIP(cusip string) (string, error) {
	cusip = strings.TrimSpace(cusip)

	cusip = doubleSpaceRegexp.ReplaceAllString(cusip, " ")

	for _, r := range []*regexp.Regexp{corpBadSpacesRegexp} {
		if r.MatchString(cusip) {
			return "", errors.CorporationInvalidCUSIP
		}
	}

	return cusip, nil
}
