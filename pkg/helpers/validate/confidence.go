package validate

import "github.com/Suburbia-io/dashboard/pkg/errors"

func Confidence(c float64) error {
	if c < 0 || c > 1 {
		return errors.InvalidConfidence
	}
	return nil
}
