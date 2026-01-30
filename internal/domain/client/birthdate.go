package client

import (
	"time"

	"github.com/edgarSucre/mye"
)

type Birthdate struct {
	val *time.Time
}

//nolint:errcheck
func NewBirthdate(t *time.Time) (Birthdate, error) {
	if t == nil {
		return Birthdate{}, nil
	}

	err := mye.New(mye.CodeInvalid, "birthdate_creation_failed", "validation failed").
		WithUserMsg("birthdate validation failed")

	if t.IsZero() {
		err.WithField("birthdate", "birthdate must be a valid date")
	}

	if t.After(time.Now()) {
		err.WithField("birthdate", "birthdate can't be set in the future")
	}

	if err.HasFields() {
		return Birthdate{}, err
	}

	return Birthdate{val: t}, nil
}

func (b Birthdate) IsValid() bool {
	if b.val == nil {
		return true
	}

	return !b.val.IsZero() && time.Now().After(*b.val)
}

func (b Birthdate) String() string {
	if b.val == nil {
		return ""
	}

	return b.val.Format(time.DateOnly)
}

func (b Birthdate) Time() *time.Time {
	return b.val
}
