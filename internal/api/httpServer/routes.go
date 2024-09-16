package httpServer

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *server) setRoutes() {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)

	mux.Route("/api", func(mux chi.Router) {
		mux.Get("/", s.home)
	})

	s.mux = mux
}
