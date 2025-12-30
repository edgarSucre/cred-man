package banks

import (
	"context"
)

type CreateBankService interface {
	Execute(context.Context, CreateBankCmd) (CreateBankResult, error)
}
