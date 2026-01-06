package client

import (
	"github.com/edgarSucre/mye"
	"github.com/google/uuid"
)

type ID struct {
	value uuid.UUID
}

func NewID(s string) (ID, error) {
	err := mye.New(mye.CodeInvalid, "client_id_creation_failed", "validation failed").
		WithUserMsg("client_id validation failed")

	if len(s) == 0 {
		return ID{}, err.WithField("client_id", "client can't be empty")
	}

	id, uuIDErr := uuid.Parse(s)
	if uuIDErr != nil {
		return ID{}, err.WithField("client_id", "client_id most be valid uuid v4")
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
