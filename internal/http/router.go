package http

import (
	"net/http"

	clientSvc "github.com/edgarSucre/crm/internal/client"
	"github.com/edgarSucre/crm/internal/http/client"
)

func addRoutes(mux *http.ServeMux, clientService *clientSvc.Service) {
	mux.Handle("GET /clients/:id", client.HandleGetClient(clientService))
	mux.Handle("POST /clients", client.HandleCreateClient(clientService))
	mux.HandleFunc("GET /health", handleHealth)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
	w.Write(nil)
}
