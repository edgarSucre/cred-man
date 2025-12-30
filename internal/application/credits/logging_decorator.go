package credits

import (
	"context"
	"log/slog"
)

/* ============================================================================================== */
/*                                     CreateCredit Decorator                                     */
/* ============================================================================================== */

type createCreditLoggerDecorator struct {
	logger *slog.Logger
	svc    CreateCreditService
}

func NewCreateCreditLoggerDecorator(
	svc CreateCreditService,
	logger *slog.Logger,
) CreateCreditService {
	return createCreditLoggerDecorator{logger, svc}
}

func (ld createCreditLoggerDecorator) Execute(
	ctx context.Context,
	cmd CreateCreditCommand,
) (CreditResult, error) {
	args := []any{
		slog.String("action", "create-credit"),
	}

	logger := ld.logger.With(args...)
	logger.Info("starting", slog.Any("payload", cmd))

	credit, err := ld.svc.Execute(ctx, cmd)
	if err != nil {
		logger.Error("result", slog.Any("failure", err))

		return credit, err
	}

	logger.Info("success")

	return credit, nil
}

/* ============================================================================================== */
/*                                   GetCredit Logger Decorator                                   */
/* ============================================================================================== */

type getCreditLoggerDecorator struct {
	logger *slog.Logger
	svc    GetCreditService
}

func NewGetCreditLoggerDecorator(
	svc GetCreditService,
	logger *slog.Logger,
) GetCreditService {
	return getCreditLoggerDecorator{logger, svc}
}

func (ld getCreditLoggerDecorator) Execute(
	ctx context.Context,
	cmd GetCreditCommand,
) (CreditResult, error) {
	args := []any{
		slog.String("action", "get-credit"),
	}

	logger := ld.logger.With(args...)
	logger.Info("starting", slog.Any("id", cmd.ID))

	resp, err := ld.svc.Execute(ctx, cmd)

	if err != nil {
		logger.Error("result", slog.Any("failure", err))
	} else {
		logger.Info("result", slog.Any("success", resp))
	}

	return resp, err
}

/* ============================================================================================== */
/*                                 ApproveCredit Logger Decorator                                 */
/* ============================================================================================== */

type approveCreditLoggerDecorator struct {
	logger *slog.Logger
}

func NewApproveCreditLoggerDecorator(logger *slog.Logger) ApproveCreditService {
	return approveCreditLoggerDecorator{logger}
}

func (ld approveCreditLoggerDecorator) Execute(ctx context.Context, creditID string) {
	ld.logger.Info("credit-approved", slog.String("creditID", creditID))
}

/* ============================================================================================== */
/*                                  RejectCredit Logger Decorator                                 */
/* ============================================================================================== */

type rejectCreditLoggerDecorator struct {
	logger *slog.Logger
}

func NewRejectCreditLoggerDecorator(logger *slog.Logger) RejectCreditService {
	return rejectCreditLoggerDecorator{logger}
}

func (ld rejectCreditLoggerDecorator) Execute(ctx context.Context, creditID string) {
	ld.logger.Info("credit-rejected", slog.String("creditID", creditID))
}

/* ============================================================================================== */
/*                                 ProcessCredit Logger Decorator                                 */
/* ============================================================================================== */

type processCreditLoggerDecorator struct {
	logger *slog.Logger
	svc    ProcessCreditService
}

func NewProcessCreditLoggerDecorator(
	svc ProcessCreditService,
	logger *slog.Logger,
) ProcessCreditService {
	return processCreditLoggerDecorator{logger, svc}
}

func (ld processCreditLoggerDecorator) Execute(
	ctx context.Context,
	cmd ProcessCreditCommand,
) error {
	args := []any{
		slog.String("action", "process-credit"),
	}

	logger := ld.logger.With(args...)
	logger.Info("processing", slog.Any("event", cmd))

	if err := ld.svc.Execute(ctx, cmd); err != nil {
		logger.Error("result", slog.Any("failure", err))
	}

	logger.Info("result", slog.String("success", "-"))

	return nil
}
