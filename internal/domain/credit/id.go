package credit

import (
	"fmt"

	"github.com/google/uuid"
)

type ID struct {
	value uuid.UUID
}

func NewIDFromString(s string) (ID, error) {
	if len(s) == 0 {
		return ID{}, ErrInvalidCreditID
	}

	id, err := uuid.Parse(s)
	if err != nil {
		return ID{}, ErrInvalidCreditID
	}

	return ID{value: id}, nil
}

func NewIDFromUUID(id uuid.UUID) (ID, error) {
	var emptyID uuid.UUID

	if id == emptyID {
		return ID{}, fmt.Errorf("credit.NewIDFromUUID > %w", ErrInvalidCreditID)
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
