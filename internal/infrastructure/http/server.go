package http

import (
	"log/slog"
	"net/http"

	"github.com/edgarSucre/crm/internal/infrastructure/config"
	"github.com/edgarSucre/crm/internal/infrastructure/thttp/httputils"
	"github.com/edgarSucre/crm/pkg/terror"
)

type ServerParams struct {
	bankHandler   BankHandler
	clientHandler ClientHandler
	creditHandler CreditHandler
	Logger        *slog.Logger
}

func NewServer(
	cfg config.Config,
	bankHandler BankHandler,
	clientHandler ClientHandler,
	creditHandler CreditHandler,
	logger *slog.Logger,
) http.Handler {
	mux := http.NewServeMux()

	addRoutes(
		mux,
		bankHandler,
		clientHandler,
		creditHandler,
	)

	var handler http.Handler = mux

	handler = httputils.RequestLogger(handler, logger)

	return handler
}

var (
	ErrNoBankHandler   = terror.Internal.New("http-bad-config", "bank handler is missing")
	ErrNoClientHandler = terror.Internal.New("http-bad-config", "client handler is missing")
	ErrNoCreditHandler = terror.Internal.New("http_bad_config", "credit handler is missing")
	ErrNoLogger        = terror.Internal.New("http-bad-config", "logger is missing")
)
