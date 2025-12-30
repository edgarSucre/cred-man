package client

import "github.com/google/uuid"

type ID struct {
	value uuid.UUID
}

func NewID(s string) (ID, error) {
	if len(s) == 0 {
		return ID{}, ErrInvalidClientID
	}

	id, err := uuid.Parse(s)
	if err != nil {
		return ID{}, ErrInvalidClientID
	}

	return ID{value: id}, nil
}

func (id ID) IsEmpty() bool {
	return len(id.value) == 0
}

func (id ID) String() string {
	return id.value.String()
}

func (id ID) UUID() uuid.UUID {
	return id.value
}
