package bank

type Bank struct {
	id       ID
	name     string
	bankType Type
}

func New(n string, t Type) (Bank, error) {
	var newBank Bank
	if len(n) == 0 {
		return newBank, ErrInvalidBankName
	}

	if !t.IsValid() {
		return newBank, ErrBankTypeInvalid
	}

	newBank.name = n
	newBank.bankType = t

	return newBank, nil
}

func (b Bank) ID() ID {
	return b.id
}

func (b Bank) Name() string {
	return b.name
}

func (b Bank) Type() Type {
	return b.bankType
}

// no validations needed
func Rehydrate(id ID, name string, t Type) Bank {
	return Bank{
		id:       id,
		name:     name,
		bankType: t,
	}
}
