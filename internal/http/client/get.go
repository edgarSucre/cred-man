package client

import (
	"net/http"

	"github.com/edgarSucre/crm/internal/client"
	"github.com/edgarSucre/crm/internal/http/utils"
	"github.com/google/uuid"
)

func HandleGetClient(svc *client.Service) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) error {
		clientId := r.PathValue("id")

		id, err := uuid.Parse(clientId)
		if err != nil {
			return err
		}

		cl, err := svc.GetClient(r.Context(), id)

		resp := new(CreateClientResponse)
		resp.FromDomain(cl)

		return utils.Marshal(w, resp)
	}

	return utils.ErrorHandlerFunc(fn)
}
