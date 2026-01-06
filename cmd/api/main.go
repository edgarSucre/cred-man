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

	"github.com/edgarSucre/crm/internal/application/banks"
	"github.com/edgarSucre/crm/internal/application/clients"
	"github.com/edgarSucre/crm/internal/application/credits"
	"github.com/edgarSucre/crm/internal/infrastructure/config"
	"github.com/edgarSucre/crm/internal/infrastructure/db/repository"
	"github.com/edgarSucre/crm/internal/infrastructure/events"
	chttp "github.com/edgarSucre/crm/internal/infrastructure/http"
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
		"HTTP_HOST":      "",
		"HTTP_PORT":      "",
		"REDIS_ADDR":     "",
	}

	cfg, err := config.LoadConfig(env)
	if err != nil {
		return err
	}

	pool, err := pgxpool.New(ctx, cfg.DbConn)
	if err != nil {
		return err
	}

	defer pool.Close()

	bankRepository := repository.NewBankRepository(pool)
	clientRepository := repository.NewClientRepository(pool)
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

	createBank, err := banks.NewCreateBankService(bankRepository)
	if err != nil {
		return err
	}

	createBank = banks.NewCreateBankLoggerDecorator(createBank, logger)

	createClient, err := clients.NewCreateClientService(clientRepository)
	if err != nil {
		return err
	}

	createClient = clients.NewCreateClientLoggerDecorator(createClient, logger)

	getClient, err := clients.NewGetClientService(clientRepository)
	if err != nil {
		return err
	}

	getClient = clients.NewGetClientLoggerDecorator(getClient, logger)

	createCredit, err := credits.NewCreateCreditService(
		bankRepository,
		clientRepository,
		creditRepository,
		eventBus,
		transactionManager,
	)
	if err != nil {
		return err
	}

	createCredit = credits.NewCreateCreditLoggerDecorator(createCredit, logger)

	getCredit := credits.NewGetCreditService(creditRepository)
	getCredit = credits.NewGetCreditLoggerDecorator(getCredit, logger)

	/* ========================================================================================== */
	/*                                        HTTP Handlers                                       */
	/* ========================================================================================== */

	clientHandler, err := chttp.NewClientHandler(createClient, getClient)
	if err != nil {
		return err
	}

	bankHandler, err := chttp.NewBankHandler(createBank)
	if err != nil {
		return err
	}

	creditHandler, err := chttp.NewCreditHandler(createCredit, getCredit)
	if err != nil {
		return err
	}

	srv, err := chttp.NewServer(
		&bankHandler,
		&clientHandler,
		&creditHandler,
		logger,
	)

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

	/* ========================================================================================== */
	/*                                      start HTTP server                                     */
	/* ========================================================================================== */

	go func() {
		logger.Info(fmt.Sprintf("http server listening on: %v", cfg.HttpPort))

		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error(fmt.Sprintf("error listening and serving: %s\n", err))
			cancel()

			return
		}

		logger.Info("server shutting down..")
	}()

	/* ========================================================================================== */
	/*                                          shutdown                                          */
	/* ========================================================================================== */

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
