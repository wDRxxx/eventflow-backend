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

		mux.Route("/events", func(mux chi.Router) {
			mux.Get("/{url-title}", s.event)
			mux.Post("/", s.createEvent)
			mux.Put("/{url-title}", s.updateEvent)
			mux.Delete("/{url-title}", s.deleteEvent)
		})

		mux.Route("/tickets", func(mux chi.Router) {
			mux.Get("/{id}", s.ticket)
		})

		mux.Route("/auth", func(mux chi.Router) {
			mux.Post("/register", s.register)
			mux.Post("/login", s.login)
			mux.Post("/refresh", s.refresh)
		})
	})

	s.mux = mux
}
