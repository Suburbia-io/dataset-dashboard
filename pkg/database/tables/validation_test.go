package tables

import (
	"github.com/Suburbia-io/dashboard/pkg/errors"
	"github.com/Suburbia-io/dashboard/pkg/helpers/crypto"
)

func invalidUUID() (string, error) {
	uuid := crypto.NewUUID()
	return uuid[1:], errors.InvalidUUID
}

func invalidUUIDPtr() (*string, error) {
	uuid := crypto.NewUUID()[1:]
	return &uuid, errors.InvalidUUID
}

func invalidEmail() (string, error) {
	return "abc.com", errors.InvalidEmail
}

func invalidHumanName() (string, error) {
	return "", errors.InvalidHumanName
}

func invalidTrimmedString() (string, error) {
	return " ", nil
}

func invalidTrimmedNullableString() (*string, error) {
	return nil, nil
}

func invalidNonEmptyString() (string, error) {
	return " \n", errors.EmptyString
}

func invalidSlug() (string, error) {
	return "abc!", errors.InvalidSlug
}

func invalidTagType() (string, error) {
	return "abc def", errors.InvalidTagType
}

func invalidPath() (string, error) {
	return "/root", errors.InvalidPath
}
