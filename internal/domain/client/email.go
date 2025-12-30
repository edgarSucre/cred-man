package client

import "net/mail"

type Email struct {
	val string
}

func NewEmail(e string) (Email, error) {
	var email Email

	if len(e) == 0 {
		return email, ErrInvalidEmail
	}

	if _, err := mail.ParseAddress(e); err != nil {
		return email, ErrInvalidEmail
	}

	email.val = e

	return email, nil
}

func (e Email) IsValid() bool {
	_, err := mail.ParseAddress(e.val)

	return len(e.val) > 0 && err == nil
}

func (e Email) String() string {
	return e.val
}
