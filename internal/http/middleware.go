package http

import (
	"log/slog"
	"net/http"
)

func logHandler(fn http.Handler, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger = logger.With(
			"method", r.Method,
			"path", r.URL.Path,
			"X-Requested-ID", r.Header.Get("X-Requested-ID"),
		)

		logger.Info("request")
	})
}
