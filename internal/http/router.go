package http

import (
	"net/http"

	"github.com/edgarSucre/crm/internal/http/client"
	"github.com/edgarSucre/crm/pkg"
)

func addRoutes(mux *http.ServeMux, clientService pkg.ClientService) {
	mux.Handle("GET /clients/:id", client.HandleGetClient(clientService))
	mux.Handle("POST /clients", client.HandleCreateClient(clientService))
	mux.HandleFunc("GET /health", handleHealth)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
	w.Write(nil)
}
