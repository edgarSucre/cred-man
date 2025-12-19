package http

import (
	"log/slog"
	"net/http"

	"github.com/edgarSucre/crm/internal/client"
	"github.com/edgarSucre/crm/internal/config"
)

func NewServer(cfg config.Config, logger *slog.Logger, clientSvc *client.Service) http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux, clientSvc)

	var handler http.Handler = mux

	handler = logHandler(handler, logger)

	return handler
}
