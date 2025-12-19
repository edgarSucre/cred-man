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

	"github.com/edgarSucre/crm/internal/client"
	"github.com/edgarSucre/crm/internal/config"
	"github.com/edgarSucre/crm/internal/db/repository"
	chttp "github.com/edgarSucre/crm/internal/http"
)

func run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	logLevel := new(slog.LevelVar)
	opts := &slog.HandlerOptions{Level: logLevel}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))

	repo, err := repository.NewRepository(ctx, cfg.DbConn)
	if err != nil {
		return fmt.Errorf("failed to connect to the database: %w", err)
	}

	clientService, err := client.NewService(repo)
	if err != nil {
		return fmt.Errorf("failed to create client-service: %w", err)
	}

	srv := chttp.NewServer(cfg, logger, clientService)
	httpServer := &http.Server{
		Addr:    net.JoinHostPort(cfg.Host, cfg.HttpPort),
		Handler: srv,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
			cancel()

			return
		}

		fmt.Fprintln(os.Stdout, " server shutting down..")
	}()

	return nil
}

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
