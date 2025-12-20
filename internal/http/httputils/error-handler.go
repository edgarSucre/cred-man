package httputils

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/edgarSucre/crm/pkg/terror"
)

type ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request) error

func (f ErrorHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := f(w, r)
	if err != nil {
		handleError(w, err)
	}
}

func handleError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	var myErr terror.Error

	if errors.As(err, &myErr) {
		unwrapped := terror.Unwrap(err)
		vErr := unwrapped.(terror.Error)

		raw, status := getErrorResponse(vErr)
		w.WriteHeader(status)

		if _, err := w.Write(raw); err != nil {
			http.Error(w, "internal error", 500)
			return
		}

		return
	}

	http.Error(w, "internal error", 500)
}

func getErrorResponse(err terror.Error) ([]byte, int) {
	if err.T == terror.Internal {
		msg := fmt.Sprintf(
			`{"code": "%s", "msg": "%s"}`,
			"server-error",
			"something went wrong, try again later",
		)

		return []byte(msg), err.HttpStatus()
	}

	return err.Raw(), err.HttpStatus()
}
