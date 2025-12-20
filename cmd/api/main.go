package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/edgarSucre/crm/internal/config"
	"github.com/edgarSucre/crm/internal/decorators"
	"github.com/edgarSucre/crm/internal/infrastructure/db/repository"
	"github.com/edgarSucre/crm/internal/infrastructure/events"
	chttp "github.com/edgarSucre/crm/internal/infrastructure/http"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func loadConfig() (config.Config, error) {
	_ = godotenv.Load()

	env := map[string]string{
		"GOOSE_DBSTRING": "",
		"HTTP_HOST":      "",
		"HTTP_PORT":      "",
		"REDIS_ADDR":     "",
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
		Host:      env["HTTP_HOST"],
		HttpPort:  env["HTTP_PORT"],
		RedisAddr: env["REDIS_ADDR"],
	}, nil
}

func run(ctx context.Context, logger *slog.Logger) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	repo, err := repository.NewRepository(ctx, cfg.DbConn)
	if err != nil {
		return err
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: "",
		DB:       0,
	})

	clientService, err := decorators.NewClientServiceWithDecorators(repo, logger)
	if err != nil {
		return err
	}

	bankService, err := decorators.NewBankServiceWithDecorators(repo, logger)
	if err != nil {
		return err

	}

	eventBus, err := events.NewStreamBus(redisClient, "domain-events")
	if err != nil {
		return err
	}

	creditService, err := decorators.NewCreditServiceWithDecorators(repo, eventBus, logger)

	srv, err := chttp.NewServer(cfg, chttp.ServerParams{
		BankService:   bankService,
		ClientService: clientService,
		CreditService: creditService,
		Logger:        logger,
	})
	if err != nil {
		return err
	}

	httpServer := &http.Server{
		Addr:    net.JoinHostPort(cfg.Host, cfg.HttpPort),
		Handler: srv,
		BaseContext: func(l net.Listener) context.Context {
			return ctx
		},
	}

	go func() {
		logger.Info(fmt.Sprintf("http server listening on: %v", cfg.HttpPort))

		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error(fmt.Sprintf("error listening and serving: %s\n", err))
			cancel()

			return
		}

		logger.Info("server shutting down..")
	}()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		<-ctx.Done()

		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 10*time.Second)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			logger.Error(fmt.Sprintf("error shutting down http server: %s", err))
		}
	}()

	wg.Wait()

	return nil
}

func main() {
	ctx := context.Background()

	logLevel := new(slog.LevelVar)
	opts := &slog.HandlerOptions{Level: logLevel}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	logger = logger.With(slog.String("api", "credit-management"))

	if err := run(ctx, logger); err != nil {
		logger.Error(err.Error())

		os.Exit(1)
	}
}
