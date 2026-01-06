package http

import (
	"net/http"
)

func addRoutes(
	mux *http.ServeMux,
	bankHandler *BankHandler,
	clientHandler *ClientHandler,
	creditHandler *CreditHandler,
) {
	mux.Handle("POST /clients", HandleCreateClient(clientHandler.createClient))
	mux.Handle("POST /banks", HandleCreateBank(bankHandler.createBank))
	mux.Handle("GET /credits/{id}", HandleGetCredit(creditHandler.getCredit))
	mux.Handle("POST /credits", HandleCreateCredit(creditHandler.createCredit))

	mux.HandleFunc("GET /health", handleHealth)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
	w.Write(nil)
}
