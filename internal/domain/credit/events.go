package credit

type CreditCreated struct {
	ID string
}

func (event CreditCreated) EventName() string {
	return "creditCreated"
}
