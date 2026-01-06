package client

import (
	"net/mail"

	"github.com/edgarSucre/mye"
)

type Email struct {
	val string
}

func NewEmail(e string) (Email, error) {
	err := mye.New(mye.CodeInvalid, "email_creation_failed", "validation failed").
		WithUserMsg("email validation failed")

	if len(e) == 0 {
		err.WithField("email", "email can't be empty")
	}

	if _, parseErr := mail.ParseAddress(e); parseErr != nil {
		err.WithField("email", "email must be a valid email address")
	}

	if err.HasFields() {
		return Email{}, err
	}

	return Email{e}, nil
}

func (e Email) IsValid() bool {
	_, err := mail.ParseAddress(e.val)

	return len(e.val) > 0 && err == nil
}

func (e Email) String() string {
	return e.val
}
