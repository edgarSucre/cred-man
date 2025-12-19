package utils

import "net/http"

type ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request) error

func (f ErrorHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := f(w, r)
	if err != nil {
		http.Error(w, "system error", 500)
	}
}
