package http

import (
	"log/slog"
	"net/http"

	"github.com/edgarSucre/crm/internal/config"
	"github.com/edgarSucre/crm/internal/infrastructure/http/httputils"
	"github.com/edgarSucre/crm/pkg/domain"
	"github.com/edgarSucre/crm/pkg/terror"
)

type ServerParams struct {
	BankService   domain.BankService
	ClientService domain.ClientService
	CreditService domain.CreditService
	Logger        *slog.Logger
}

func NewServer(cfg config.Config, params ServerParams) (http.Handler, error) {
	if err := params.Validate(); err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	addRoutes(mux, params.BankService, params.ClientService, params.CreditService)

	var handler http.Handler = mux

	handler = httputils.RequestLogger(handler, params.Logger)

	return handler, nil
}

var (
	ErrNoBankService   = terror.Internal.New("http-bad-config", "bank service is missing")
	ErrNoClientService = terror.Internal.New("http-bad-config", "client service is missing")
	ErrNoLogger        = terror.Internal.New("http-bad-config", "logger is missing")
)

func (params ServerParams) Validate() error {
	if params.BankService == nil {
		return ErrNoBankService
	}

	if params.ClientService == nil {
		return ErrNoClientService
	}

	if params.Logger == nil {
		return ErrNoLogger
	}

	return nil
}
