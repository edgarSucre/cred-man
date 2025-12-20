package httputils

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
)

func RequestLogger(next http.Handler, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attrs := []any{
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("X-Requested-ID", r.Header.Get("X-Requested-ID")),
		}

		lw := &LoggerWriter{ResponseWriter: w}

		next.ServeHTTP(lw, r)

		attrs = append(attrs, slog.Int("status", lw.status))

		logger.Info("request", attrs...)
	})
}

type LoggerWriter struct {
	http.ResponseWriter
	buf    bytes.Buffer
	status int
}

func (lw *LoggerWriter) Write(data []byte) (int, error) {
	writer := io.MultiWriter(lw.ResponseWriter, &lw.buf)

	return writer.Write(data)
}

func (lw *LoggerWriter) WriteHeader(status int) {
	lw.status = status
	lw.ResponseWriter.WriteHeader(status)
}
