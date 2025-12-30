package bank

type Name struct {
	value string
}

func NewName(s string) (Name, error) {
	if len(s) == 0 {
		return Name{}, ErrInvalidBankName
	}

	return Name{value: s}, nil
}

func (n Name) IsEmpty() bool {
	return len(n.value) == 0
}

func (n Name) String() string {
	return n.value
}
