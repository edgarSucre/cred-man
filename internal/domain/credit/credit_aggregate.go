package credit

type CreditAggregate struct {
	clientCredits []Credit
	creditType    CreditType
	id            ID
	status        CreditStatus
}

func RehydrateAggregate(
	clientCredits []Credit,
	creditType CreditType,
	id ID,
	status CreditStatus,
) *CreditAggregate {
	return &CreditAggregate{clientCredits, creditType, id, status}
}

func (ca *CreditAggregate) Status() CreditStatus {
	return ca.status
}

func (ca *CreditAggregate) ID() ID {
	return ca.id
}

func (ca *CreditAggregate) Process() {
	mortgages := make([]Credit, 0, len(ca.clientCredits))
	autos := make([]Credit, 0, len(ca.clientCredits))
	commercial := make([]Credit, 0, len(ca.clientCredits))

	for _, v := range ca.clientCredits {
		switch v.CreditType() {
		case CreditTypeMortgage:
			mortgages = append(mortgages, v)
		case CreditTypeAuto:
			autos = append(autos, v)
		default:
			commercial = append(commercial, v)
		}
	}

	if ca.creditType == CreditTypeAuto && len(autos) >= 2 {
		ca.status = CreditStatusRejected
		return
	}

	if ca.creditType == CreditTypeCommercial && len(commercial) >= 3 {
		ca.status = CreditStatusRejected
		return
	}

	if ca.creditType == CreditTypeMortgage && len(mortgages) >= 4 {
		ca.status = CreditStatusRejected
		return
	}

	ca.status = CreditStatusApproved
}
