package client

import "strings"

type FullName struct {
	val string
}

func NewFullName(n string) (FullName, error) {
	var fullName FullName

	if len(strings.TrimSpace(n)) == 0 {
		return fullName, ErrInvalidClientFullName
	}

	fullName.val = n

	return fullName, nil
}

func (fn FullName) IsEmpty() bool {
	return len(fn.val) == 0
}

func (fn FullName) String() string {
	return fn.String()
}
