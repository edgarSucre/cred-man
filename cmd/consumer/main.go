package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	"github.com/edgarSucre/crm/internal/config"
	"github.com/edgarSucre/crm/internal/decorators"
	"github.com/edgarSucre/crm/internal/handlers"
	"github.com/edgarSucre/crm/internal/infrastructure/db/repository"
	"github.com/edgarSucre/crm/internal/infrastructure/events"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func loadConfig() (config.Config, error) {
	_ = godotenv.Load()

	env := map[string]string{
		"GOOSE_DBSTRING": "",
		"REDIS_ADDR":     "",
		"CONSUMER":       "",
	}

	for key := range env {
		val := os.Getenv(key)

		if len(val) == 0 {
			return config.Config{}, config.ErrLoadConfig(key)
		}

		env[key] = val
	}

	return config.Config{
		DbConn:    env["GOOSE_DBSTRING"],
		RedisAddr: env["REDIS_ADDR"],
		Consumer:  env["CONSUMER"],
	}, nil
}

func run(ctx context.Context, logger *slog.Logger) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	logger = logger.With(slog.String("consumer", cfg.Consumer))

	repo, err := repository.NewRepository(ctx, cfg.DbConn)
	if err != nil {
		return err
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: "",
		DB:       0,
	})

	eventBus, err := events.NewStreamBus(redisClient, "domain-events")
	if err != nil {
		return err
	}

	creditService, err := decorators.NewCreditServiceWithDecorators(repo, eventBus, logger)
	if err != nil {
		return err
	}

	creditHandlers, err := handlers.GetCreditHandlers(creditService)
	if err != nil {
		return err
	}

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
