package client

import "time"

type Birthdate struct {
	val *time.Time
}

func NewBirthdate(t *time.Time) (Birthdate, error) {
	if t == nil {
		return Birthdate{}, nil
	}

	if t.IsZero() {
		return Birthdate{}, ErrInvalidBirthdate
	}

	if t.After(time.Now()) {
		return Birthdate{}, ErrInvalidBirthdate
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
