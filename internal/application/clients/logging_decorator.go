package clients

import (
	"context"
	"log/slog"
)

type CreateClientLoggerDecorator struct {
	svc    CreateClientService
	logger *slog.Logger
}

func NewCreateClientLoggerDecorator(
	svc CreateClientService,
	logger *slog.Logger,
) CreateClientService {
	return CreateClientLoggerDecorator{svc: svc, logger: logger}
}

func (ld CreateClientLoggerDecorator) CreateClient(
	ctx context.Context,
	cmd CreateClientCommand,
) (ClientResult, error) {
	args := []any{
		slog.String("action", "create-client"),
	}

	logger := ld.logger.With(args...)
	logger.Info("starting", slog.Any("payload", cmd))

	resp, err := ld.svc.CreateClient(ctx, cmd)

	if err != nil {
		logger.Error("create-client-failure", slog.Any("cause", err))
	} else {
		logger.Info("create-client-success", slog.Any("result", resp))
	}

	return resp, err
}

type GetClientLoggerDecorator struct {
	svc    GetClientService
	logger *slog.Logger
}

func NewGetClientLoggerDecorator(
	svc GetClientService,
	logger *slog.Logger,
) GetClientService {
	return GetClientLoggerDecorator{svc: svc, logger: logger}
}

func (ld GetClientLoggerDecorator) GetClient(
	ctx context.Context,
	cmd GetClientCommand,
) (ClientResult, error) {
	args := []any{
		slog.String("action", "get-client"),
	}

	logger := ld.logger.With(args...)
	logger.Info("starting", slog.Any("payload", cmd))

	resp, err := ld.svc.GetClient(ctx, cmd)

	if err != nil {
		logger.Error("get-client-failure", slog.Any("cause", err))
	} else {
		logger.Info("get-client-success", slog.Any("result", resp))
	}

	return resp, err
}
