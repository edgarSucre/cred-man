package bank

import (
	"github.com/edgarSucre/mye"
	"github.com/google/uuid"
)

type ID struct {
	value uuid.UUID
}

func NewID(s string) (ID, error) {
	err := mye.New(mye.CodeInvalid, "bank_id_error", "failed to create bankID").
		WithUserMsg("bank id fails validation")

	if len(s) == 0 {
		return ID{}, err.WithField("bank_id", "bank_id can't be empty")
	}

	id, uuidErr := uuid.Parse(s)
	if uuidErr != nil {
		return ID{}, err.WithField("bank_id", "bank_id must be a valid uuid v4")
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
