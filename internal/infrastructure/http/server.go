package http

import (
	"log/slog"
	"net/http"

	"github.com/edgarSucre/crm/internal/infrastructure/http/httputils"
	"github.com/edgarSucre/mye"
)

type ServerParams struct {
	bankHandler   *BankHandler
	clientHandler *ClientHandler
	creditHandler *CreditHandler
	Logger        *slog.Logger
}

func NewServer(
	bankHandler *BankHandler,
	clientHandler *ClientHandler,
	creditHandler *CreditHandler,
	logger *slog.Logger,
) (http.Handler, error) {
	if err := validateServer(bankHandler, clientHandler, creditHandler, logger); err != nil {
		return nil, err
	}

	mux := http.NewServeMux()

	addRoutes(
		mux,
		bankHandler,
		clientHandler,
		creditHandler,
	)

	var handler http.Handler = mux

	handler = httputils.RequestLogger(handler, logger)

	return handler, nil
}

func validateServer(
	bankHandler *BankHandler,
	clientHandler *ClientHandler,
	creditHandler *CreditHandler,
	logger *slog.Logger,
) error {
	err := mye.New(
		mye.CodeInternal,
		"http_server_config_error",
		"http server params validation error",
	)

	if bankHandler == nil {
		err.WithField("bankHandler", "bank handler is missing")
	}

	if clientHandler == nil {
		err.WithField("clientHandler", "client handler is missing")
	}

	if creditHandler == nil {
		err.WithField("creditHandler", "credit handler is missing")
	}

	if logger == nil {
		err.WithField("logger", "logger is missing")
	}

	if err.HasFields() {
		return err
	}

	return nil
}
