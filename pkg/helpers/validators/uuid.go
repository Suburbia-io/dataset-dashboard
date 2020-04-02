package validators

import (
	"regexp"

	"github.com/Suburbia-io/dashboard/pkg/errors"
)

var UUIDPattern = regexp.MustCompile(
	`^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$`)

func UUID(uuid string) (string, error) {
	if UUIDPattern.MatchString(uuid) {
		return uuid, nil
	}
	return uuid, errors.InvalidUUID
}

func NullableUUID(uuid *string) (*string, error) {
	if uuid == nil {
		return nil, nil
	}
	u, err := UUID(*uuid)
	return &u, err
}

func UUIDDeleteMe(uuid string) error {
	if UUIDPattern.MatchString(uuid) {
		return nil
	}
	return errors.InvalidUUID
}
