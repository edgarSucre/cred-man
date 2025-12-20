package decorators

import (
	"context"
	"log/slog"

	"github.com/edgarSucre/crm/internal/infrastructure/db/repository"
	"github.com/edgarSucre/crm/internal/usecases/client"
	"github.com/edgarSucre/crm/pkg/domain"
	"github.com/google/uuid"
)

type clientLoggerDecorator struct {
	svc    domain.ClientService
	logger *slog.Logger
}

func (cl *clientLoggerDecorator) CreateClient(
	ctx context.Context,
	params domain.CreateClientParams,
) (*domain.Client, error) {
	args := []any{
		slog.String("action", "create-client"),
	}

	logger := cl.logger.With(args...)
	logger.Info("starting..")

	resp, err := cl.svc.CreateClient(ctx, params)

	if err != nil {
		logger.Error("create-client-failure", slog.Any("cause", err))
	} else {
		logger.Info("create-client-success", slog.Any("result", resp))
	}

	return resp, err
}

func (cl *clientLoggerDecorator) GetClient(ctx context.Context, id uuid.UUID) (*domain.Client, error) {
	return cl.svc.GetClient(ctx, id)
}

func NewClientServiceWithDecorators(
	repo repository.Querier,
	logger *slog.Logger,
) (domain.ClientService, error) {
	svc, err := client.NewService(repo)

	if err != nil {
		return nil, err
	}

	logger = logger.With("service", "client")

	return &clientLoggerDecorator{
		svc:    svc,
		logger: logger,
	}, nil
}
