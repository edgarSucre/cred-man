package credit

import "github.com/edgarSucre/mye"

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
	err := mye.New(mye.CodeInvalid, "credit_type_creation_failed", "validation failed").
		WithUserMsg("credit type validation failed")

	if len(s) == 0 {
		return CreditTypeInvalid, err.WithField("credit_type", "credit_type can't be empty")
	}

	switch s {
	case CreditTypeAuto.value:
		return CreditTypeAuto, nil
	case CreditTypeCommercial.value:
		return CreditTypeCommercial, nil
	case CreditTypeMortgage.value:
		return CreditTypeMortgage, nil
	default:
		return CreditTypeInvalid, err.WithField("credit_type", "credit_type is not a valid credit type")
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
