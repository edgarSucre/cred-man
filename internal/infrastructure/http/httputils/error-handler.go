package httputils

import (
	"net/http"

	"github.com/edgarSucre/mye/httpconv"
)

type ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request) error

func (f ErrorHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := f(w, r)

	if err != nil {
		httpconv.WriteHTTPError(w, err)
	}
}
