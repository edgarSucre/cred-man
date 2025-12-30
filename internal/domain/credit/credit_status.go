package credit

type (
	CreditStatus struct {
		value string
	}
)

var (
	CreditStatusInvalid  = CreditStatus{}
	CreditStatusPending  = CreditStatus{"pending"}
	CreditStatusApproved = CreditStatus{"approved"}
	CreditStatusRejected = CreditStatus{"rejected"}
)

func CreditStatusFromString(s string) (CreditStatus, error) {
	if len(s) == 0 {
		return CreditStatusInvalid, ErrInvalidCreditStatus
	}

	switch s {
	case CreditStatusPending.value:
		return CreditStatusPending, nil
	case CreditStatusApproved.value:
		return CreditStatusApproved, nil
	case CreditStatusRejected.value:
		return CreditStatusRejected, nil
	default:
		return CreditStatusInvalid, ErrInvalidCreditStatus
	}
}

func (ct CreditStatus) IsInValid() bool {
	return ct == CreditStatusInvalid
}

func (ct CreditStatus) String() string {
	return ct.value
}
