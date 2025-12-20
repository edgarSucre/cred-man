package decorators

import (
	"context"
	"log/slog"

	"github.com/edgarSucre/crm/internal/infrastructure/db/repository"
	"github.com/edgarSucre/crm/internal/usecases/bank"
	"github.com/edgarSucre/crm/pkg/domain"
)

type bankLoggerDecorator struct {
	svc    domain.BankService
	logger *slog.Logger
}

func (bl *bankLoggerDecorator) CreateBank(
	ctx context.Context,
	params domain.CreateBankParams,
) (*domain.Bank, error) {
	args := []any{
		slog.String("action", "create-bank"),
	}

	logger := bl.logger.With(args...)
	logger.Info("starting", slog.Any("payload", params))

	resp, err := bl.svc.CreateBank(ctx, params)

	if err != nil {
		logger.Error("result", slog.Any("failure", err))
	} else {
		logger.Info("result", slog.Any("success", resp))
	}

	return resp, err
}

func NewBankServiceWithDecorators(
	repo repository.Querier,
	logger *slog.Logger,
) (domain.BankService, error) {
	svc, err := bank.NewService(repo)

	if err != nil {
		return nil, err
	}

	logger = logger.With("service", "bank")

	return &bankLoggerDecorator{
		svc:    svc,
		logger: logger,
	}, nil
}
