package decorators

import (
	"context"
	"log/slog"

	"github.com/edgarSucre/crm/internal/infrastructure/db/repository"
	"github.com/edgarSucre/crm/internal/usecases/credit"
	"github.com/edgarSucre/crm/pkg/domain"
	"github.com/edgarSucre/crm/pkg/events"
	"github.com/google/uuid"
)

type creditLoggerDecorator struct {
	svc    domain.CreditService
	logger *slog.Logger
}

func (cl *creditLoggerDecorator) CreateCredit(
	ctx context.Context,
	params domain.CreateCreditParams,
) error {
	args := []any{
		slog.String("action", "create-credit"),
	}

	logger := cl.logger.With(args...)
	logger.Info("starting", slog.Any("payload", params))

	if err := cl.svc.CreateCredit(ctx, params); err != nil {
		logger.Error("result", slog.Any("failure", err))

		return err
	}

	logger.Info("success")

	return nil
}

func (cl *creditLoggerDecorator) ProcessCredit(
	ctx context.Context,
	event domain.CreditCreated,
) error {
	args := []any{
		slog.String("action", "process-credit"),
	}

	logger := cl.logger.With(args...)
	logger.Info("starting", slog.Any("event", event))

	if err := cl.svc.ProcessCredit(ctx, event); err != nil {
		logger.Error("result", slog.Any("failure", err))
	}

	logger.Info("success")

	return nil
}

func (cl *creditLoggerDecorator) GetCredit(
	ctx context.Context,
	id uuid.UUID,
) (*domain.Credit, error) {
	args := []any{
		slog.String("action", "get-credit"),
	}

	logger := cl.logger.With(args...)
	logger.Info("starting", slog.Any("id", id))

	resp, err := cl.svc.GetCredit(ctx, id)

	if err != nil {
		logger.Error("result", slog.Any("failure", err))
	} else {
		logger.Info("result", slog.Any("success", resp))
	}

	return resp, err
}

func (cl *creditLoggerDecorator) ProcessApproval(
	ctx context.Context,
	event domain.CreditApproved,
) error {
	args := []any{
		slog.String("action", "process-credit-approval"),
	}

	logger := cl.logger.With(args...)
	logger.Info("processing", slog.Any("event", event))

	return nil
}

func (cl *creditLoggerDecorator) ProcessRejection(
	ctx context.Context,
	event domain.CreditRejected,
) error {
	args := []any{
		slog.String("action", "process-credit-rejected"),
	}

	logger := cl.logger.With(args...)
	logger.Info("processing", slog.Any("event", event))

	return nil
}

func NewCreditServiceWithDecorators(
	repo *repository.Queries,
	bus events.Bus,
	logger *slog.Logger,
) (domain.CreditService, error) {
	svc, err := credit.NewService(repo, bus)

	if err != nil {
		return nil, err
	}

	logger = logger.With("service", "credit")

	return &creditLoggerDecorator{
		svc:    svc,
		logger: logger,
	}, nil
}
