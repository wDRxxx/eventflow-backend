package httpServer

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/pkg/errors"
	"github.com/rs/cors"

	"github.com/wDRxxx/eventflow-backend/internal/metrics"
)

var errInvalidAuthHeader = errors.New("invalid auth header")

func (s *server) enableCORS(next http.Handler) http.Handler {
	return cors.New(cors.Options{
		AllowCredentials: true,
		AllowedOrigins:   s.httpConfig.Origins(),
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Authorization"},
	}).Handler(next)
}

func (s *server) authRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Add("Vary", "Authorization")

		_, _, err := s.getAndVerifyHeaderToken(r)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *server) metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metrics.IncRequestCounter()

		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)
		respTime := time.Since(start)

		status := ww.Status()
		uri := r.RequestURI

		if strings.Contains(uri, "events/") {
			exploded := strings.Split(uri, "/")
			uri = strings.Join(exploded[:len(exploded)-1], "/")
		}

		if strings.Contains(uri, "?") {
			exploded := strings.Split(uri, "?")
			uri = strings.Join(exploded[:len(exploded)-1], "")
		}

		if !strings.Contains(uri, "static") {
			metrics.IncResponseCounter(status, uri)
			metrics.HistogramsResponseTimeObserve(status, uri, respTime.Seconds())
		}
	})
}
