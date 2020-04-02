package validate

import (
	"regexp"

	"github.com/Suburbia-io/dashboard/pkg/errors"
)

var uuidPattern = regexp.MustCompile(
	`^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$`)

func UUID(uuid string) error {
	if uuidPattern.MatchString(uuid) {
		return nil
	}
	return errors.InvalidUUID
}

func UUIDPtr(uuid *string) error {
	if uuid == nil {
		return nil
	}
	return UUID(*uuid)
}
