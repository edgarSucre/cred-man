package banks

import (
	"context"

	"github.com/edgarSucre/crm/internal/domain/bank"
	"github.com/edgarSucre/mye"
)

type createBank struct {
	repo Repository
}

func NewCreateBankService(repo Repository) (CreateBankService, error) {
	if repo == nil {
		return nil, mye.New(
			mye.CodeInvalid,
			"createBank_service_config_error",
			"client repository is missing",
		)
	}

	return createBank{repo: repo}, nil
}

func (svc createBank) Execute(
	ctx context.Context,
	cmd CreateBankCmd,
) (CreateBankResult, error) {
	if err := cmd.Validate(); err != nil {
		return CreateBankResult{}, err
	}

	bankType, err := bank.TypeFromString(cmd.Type)
	if err != nil {
		return CreateBankResult{}, err
	}

	dBank, err := bank.New(cmd.Name, bankType)
	if err != nil {
		return CreateBankResult{}, err
	}

	dBank, err = svc.repo.CreateBank(ctx, dBank)

	if err != nil {
		return CreateBankResult{}, err
	}

	return mapBankResult(dBank), nil
}

type CreateBankCmd struct {
	Name string
	Type string
}

func (cmd CreateBankCmd) Validate() error {
	err := mye.New(mye.CodeInvalid, "bank_creation_failed", "validation failed").
		WithUserMsg("bank validation failed")

	if len(cmd.Name) == 0 {
		err.WithField("name", "bank name is required")
	}

	if len(cmd.Type) == 0 {
		err.WithField("type", "bank type is required")
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
