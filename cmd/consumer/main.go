package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	"github.com/edgarSucre/crm/internal/application/credits"
	"github.com/edgarSucre/crm/internal/handlers"
	"github.com/edgarSucre/crm/internal/infrastructure/config"
	"github.com/edgarSucre/crm/internal/infrastructure/db/repository"
	"github.com/edgarSucre/crm/internal/infrastructure/events"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func run(ctx context.Context, logger *slog.Logger) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	/* ========================================================================================== */
	/*                                       infrastructure                                       */
	/* ========================================================================================== */

	env := map[string]string{
		"GOOSE_DBSTRING": "",
		"REDIS_ADDR":     "",
		"CONSUMER":       "",
	}

	cfg, err := config.LoadConfig(env)
	if err != nil {
		return err
	}

	logger = logger.With(slog.String("consumer", cfg.Consumer))

	pool, err := pgxpool.New(ctx, cfg.DbConn)
	if err != nil {
		return err
	}

	defer pool.Close()

	creditRepository := repository.NewCreditRepository(pool)
	transactionManager := repository.NewTransactionManager(pool)

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: "",
		DB:       0,
	})

	eventBus, err := events.NewStreamBus(redisClient, "domain-events")
	if err != nil {
		return err
	}

	/* ========================================================================================== */
	/*                                          use case                                          */
	/* ========================================================================================== */
	approveCredit := credits.NewApproveCreditLoggerDecorator(logger)
	rejectCredit := credits.NewRejectCreditLoggerDecorator(logger)

	processCredit := credits.NewProcessCreditService(eventBus, creditRepository, transactionManager)
	processCredit = credits.NewProcessCreditLoggerDecorator(processCredit, logger)

	/* ========================================================================================== */
	/*                                    domain event handlers                                   */
	/* ========================================================================================== */

	creditHandlers := handlers.GetCreditHandlers(approveCredit, processCredit, rejectCredit)

	handlers := make(map[string]events.EventHandler, 3)

	for _, h := range creditHandlers {
		handlers[h.EventName()] = h
	}

	streamConsumer, err := events.NewConsumer(events.ConsumerParams{
		Client:   redisClient,
		Consumer: cfg.Consumer,
		Stream:   "domain-events",
		Group:    "credit-management",
		Handlers: handlers,
	})

	if err != nil {
		return err
	}

	/* ========================================================================================== */
	/*                                    start event consumer                                    */
	/* ========================================================================================== */

	logger.Info("stream consumer is listening")

	if err := streamConsumer.Start(ctx); err != nil {
		logger.Error(fmt.Sprintf("error starting stream consumer: %s\n", err))
		cancel()

		return err
	}

	logger.Info("stream consumer is closed, shutting down..")

	return nil
}

func main() {
	ctx := context.Background()

	logLevel := new(slog.LevelVar)
	opts := &slog.HandlerOptions{Level: logLevel}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))

	if err := run(ctx, logger); err != nil {
		logger.Error(err.Error())

		os.Exit(1)
	}
}
