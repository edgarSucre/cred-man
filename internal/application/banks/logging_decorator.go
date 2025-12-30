package banks

import (
	"context"
	"log/slog"
)

type CreateBankLoggerDecorator struct {
	svc    CreateBankService
	logger *slog.Logger
}

func NewCreateBankLoggerDecorator(svc CreateBankService, logger *slog.Logger) CreateBankService {
	return &CreateBankLoggerDecorator{
		svc:    svc,
		logger: logger,
	}
}

func (bl *CreateBankLoggerDecorator) Execute(
	ctx context.Context,
	cmd CreateBankCmd,
) (CreateBankResult, error) {
	args := []any{
		slog.String("action", "create-bank"),
	}

	logger := bl.logger.With(args...)
	logger.Info("starting", slog.Any("payload", cmd))

	resp, err := bl.svc.Execute(ctx, cmd)

	if err != nil {
		logger.Error("result", slog.Any("failure", err))
	} else {
		logger.Info("result", slog.Any("success", resp))
	}

	return resp, err
}
