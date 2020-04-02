package sanitize

import (
	"net/mail"
	"strings"
)

func Email(email string) string {
	email = strings.ToLower(email)
	address, err := mail.ParseAddress(email)
	if err == nil {
		return address.Address
	}
	return email
}
