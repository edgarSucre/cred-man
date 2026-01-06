package credit

import (
	"github.com/edgarSucre/mye"
	"github.com/google/uuid"
)

type ID struct {
	value uuid.UUID
}

func NewIDFromString(s string) (ID, error) {
	err := mye.New(mye.CodeInvalid, "credit_id_creation_failed", "validation failed").
		WithUserMsg("credit_id creation validation failed")

	if len(s) == 0 {
		return ID{}, err.WithField("credit_id", "credit_id can't be empty")
	}

	id, uuidErr := uuid.Parse(s)
	if uuidErr != nil {
		return ID{}, err.WithField("credit_id", "credit_id is not a valid uuid v4")
	}

	return ID{value: id}, nil
}

func NewIDFromUUID(id uuid.UUID) (ID, error) {
	var emptyID uuid.UUID

	err := mye.New(mye.CodeInvalid, "credit_id_creation_failed", "validation failed").
		WithUserMsg("credit_id creation validation failed")

	if id == emptyID {
		return ID{}, err.WithField("credit_id", "credit_id can't be empty")
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

func (id ID) IsEqual(oID ID) bool {
	return id.String() == oID.String()
}
