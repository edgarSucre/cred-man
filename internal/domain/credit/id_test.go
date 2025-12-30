package credit

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewIDFromUUID(t *testing.T) {
	var uuID1 uuid.UUID

	id, err := NewIDFromUUID(uuID1)

	assert.ErrorIs(t, err, ErrInvalidCreditID)
	assert.Empty(t, id)
}
