package bank

import "github.com/edgarSucre/mye"

type Type struct {
	value string
}

func TypeFromString(t string) (Type, error) {
	err := mye.New(mye.CodeInvalid, "bank_type_creation_failed", "validation failed").
		WithUserMsg("bank type creation failed due to invalid input")

	if len(t) == 0 {
		return Type{}, err.WithField("type", "bank_type can't be empty")
	}

	switch t {
	case BankTypeGovernment.value:
		return BankTypeGovernment, nil
	case BankTypePrivate.value:
		return BankTypePrivate, nil
	default:
		return BankTypeInvalid, err.WithField("type", "type is no a valid bank type")
	}
}

var (
	BankTypeInvalid    = Type{}
	BankTypePrivate    = Type{value: "private"}
	BankTypeGovernment = Type{value: "government"}
)

func (bt Type) IsValid() bool {
	return bt == BankTypeGovernment || bt == BankTypePrivate
}

func (bt Type) String() string {
	return bt.value
}
