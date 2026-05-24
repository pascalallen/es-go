package email

import (
	"fmt"
	"net/mail"
)

type EmailAddress string

func New(s string) (EmailAddress, error) {
	if s == "" {
		return "", fmt.Errorf("email address cannot be empty")
	}
	addr, err := mail.ParseAddress(s)
	if err != nil {
		return "", fmt.Errorf("invalid email address %q: %w", s, err)
	}
	return EmailAddress(addr.Address), nil
}

func (e EmailAddress) String() string {
	return string(e)
}
