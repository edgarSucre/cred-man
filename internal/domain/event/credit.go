package event

type Event interface {
	EventName() string
}

type (
	CreditCreated struct {
		BankID   string
		ClientID string
		CreditID string
	}

	CreditApproved struct {
		CreditID string
	}

	CreditRejected struct {
		CreditID string
	}
)

func (CreditCreated) EventName() string {
	return "credit.created"
}

func (CreditApproved) EventName() string {
	return "credit.approved"
}

func (CreditRejected) EventName() string {
	return "credit.rejected"
}
