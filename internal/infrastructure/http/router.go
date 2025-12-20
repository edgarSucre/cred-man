package http

import (
	"net/http"

	"github.com/edgarSucre/crm/internal/infrastructure/http/bank"
	"github.com/edgarSucre/crm/internal/infrastructure/http/client"
	"github.com/edgarSucre/crm/pkg/domain"
)

func addRoutes(
	mux *http.ServeMux,
	bankService domain.BankService,
	clientService domain.ClientService,
) {
	// mux.Handle("GET /clients/:id", client.HandleGetClient(clientService))
	mux.Handle("POST /clients", client.HandleCreateClient(clientService))
	mux.Handle("POST /banks", bank.HandleCreateBank(bankService))

	mux.HandleFunc("GET /health", handleHealth)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
	w.Write(nil)
}
