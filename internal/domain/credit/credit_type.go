package credit

type CreditType struct {
	value string
}

var (
	CreditTypeInvalid    = CreditType{}
	CreditTypeAuto       = CreditType{"auto"}
	CreditTypeMortgage   = CreditType{"mortgage"}
	CreditTypeCommercial = CreditType{"commercial"}
)

func CreditTypeFromString(s string) (CreditType, error) {
	if len(s) == 0 {
		return CreditTypeInvalid, ErrInvalidCreditType
	}

	switch s {
	case CreditTypeAuto.value:
		return CreditTypeAuto, nil
	case CreditTypeCommercial.value:
		return CreditTypeCommercial, nil
	case CreditTypeMortgage.value:
		return CreditTypeMortgage, nil
	default:
		return CreditTypeInvalid, ErrInvalidCreditType
	}
}

func (ct CreditType) IsInvalid() bool {
	return ct == CreditTypeInvalid
}

func (ct CreditType) String() string {
	return ct.value
}

func (ct CreditType) IsEqual(ot CreditType) bool {
	return ct == ot
}
