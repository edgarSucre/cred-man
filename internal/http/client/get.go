package client

import (
	"net/http"

	"github.com/edgarSucre/crm/internal/http/httputils"
	"github.com/edgarSucre/crm/pkg"
	"github.com/google/uuid"
)

func HandleGetClient(svc pkg.ClientService) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) error {
		clientId := r.PathValue("id")

		id, err := uuid.Parse(clientId)
		if err != nil {
			return err
		}

		cl, err := svc.GetClient(r.Context(), id)

		resp := new(CreateClientResponse)
		resp.FromDomain(cl)

		return httputils.Marshal(w, resp)
	}

	return httputils.ErrorHandlerFunc(fn)
}
