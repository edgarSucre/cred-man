package credits

import "context"

type CreateCreditService interface {
	Execute(context.Context, CreateCreditCommand) (CreditResult, error)
}

type GetCreditService interface {
	Execute(context.Context, GetCreditCommand) (CreditResult, error)
}

type ProcessCreditService interface {
	Execute(context.Context, ProcessCreditCommand) error
}

type ApproveCreditService interface {
	Execute(context.Context, string)
}

type RejectCreditService interface {
	Execute(context.Context, string)
}
