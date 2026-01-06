package credit

import "github.com/edgarSucre/mye"

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
	err := mye.New(mye.CodeInvalid, "credit_status_creation_failed", "validation failed").
		WithUserMsg("credit status validation failed")

	if len(s) == 0 {
		return CreditStatusInvalid, err.WithField("credit_status", "can't be empty")
	}

	switch s {
	case CreditStatusPending.value:
		return CreditStatusPending, nil
	case CreditStatusApproved.value:
		return CreditStatusApproved, nil
	case CreditStatusRejected.value:
		return CreditStatusRejected, nil
	default:
		return CreditStatusInvalid, err.WithField("credit_status", "is not a valid credit status")
	}
}

func (ct CreditStatus) IsInValid() bool {
	return ct == CreditStatusInvalid
}

func (ct CreditStatus) String() string {
	return ct.value
}
