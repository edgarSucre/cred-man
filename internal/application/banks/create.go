package banks

import (
	"context"
	"fmt"

	"github.com/edgarSucre/crm/internal/domain/bank"
)

type createBank struct {
	repo Repository
}

func NewCreateBankService(repo Repository) CreateBankService {
	return createBank{repo: repo}
}

func (svc createBank) Execute(
	ctx context.Context,
	cmd CreateBankCmd,
) (CreateBankResult, error) {
	if err := cmd.Validate(); err != nil {
		return CreateBankResult{}, fmt.Errorf("create bank validation-error: %w", err)
	}

	bankType, err := bank.TypeFromString(cmd.Type)
	if err != nil {
		return CreateBankResult{}, fmt.Errorf("bank.TypeFromString > %w", err)
	}

	dBank, err := bank.New(cmd.Name, bankType)
	if err != nil {
		return CreateBankResult{}, fmt.Errorf("bank.New() error: %w", err)
	}

	dBank, err = svc.repo.CreateBank(ctx, dBank)

	if err != nil {
		return CreateBankResult{}, fmt.Errorf("repo.Create() error: %w", err)
	}

	return mapBankResult(dBank), nil
}

type CreateBankCmd struct {
	Name string
	Type string
}

// TODO: replace error with application errors
func (cmd CreateBankCmd) Validate() error {
	if len(cmd.Name) == 0 {
		return ErrBankNameRequired
	}

	if len(cmd.Type) == 0 {
		return ErrBankTypeInvalid
	}

	return nil
}

type CreateBankResult struct {
	ID   string
	Name string
	Type string
}

func mapBankResult(db bank.Bank) CreateBankResult {
	return CreateBankResult{
		ID:   db.ID().String(),
		Name: db.Name(),
		Type: db.Type().String(),
	}
}
