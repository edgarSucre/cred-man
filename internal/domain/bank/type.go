package bank

type Type struct {
	value string
}

func TypeFromString(t string) (Type, error) {
	if len(t) == 0 {
		return Type{}, ErrBankTypeInvalid
	}

	switch t {
	case BankTypeGovernment.value:
		return BankTypeGovernment, nil
	case BankTypePrivate.value:
		return BankTypePrivate, nil
	default:
		return BankTypeInvalid, ErrBankTypeInvalid
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
