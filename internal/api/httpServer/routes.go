package httpServer

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *server) setRoutes() {
	mux := chi.NewRouter()

	mux.Use(s.metrics)
	mux.Use(middleware.Recoverer)
	mux.Use(s.enableCORS)

	fs := http.FileServer(http.Dir(s.httpConfig.StaticDir()))
	mux.Handle("/api/static/*", http.StripPrefix("/api/static/", fs))

	mux.Route("/api", func(mux chi.Router) {
		mux.Route("/events", func(mux chi.Router) {
			mux.Get("/", s.events)
			mux.Get("/{url-title}", s.event)

			mux.Group(func(mux chi.Router) {
				mux.Use(s.authRequired)

				mux.Post("/", s.createEvent)
				mux.Put("/{url-title}", s.updateEvent)
				mux.Delete("/{url-title}", s.deleteEvent)
			})
		})

		mux.Route("/tickets", func(mux chi.Router) {
			mux.Use(s.authRequired)

			mux.Post("/", s.buyTicket)
		})

		mux.Route("/auth", func(mux chi.Router) {
			mux.Post("/register", s.register)
			mux.Post("/login", s.login)
			mux.Post("/refresh", s.refresh)
			mux.Post("/logout", s.logout)
		})

		mux.Route("/user", func(mux chi.Router) {
			mux.Use(s.authRequired)

			mux.Get("/tickets", s.userTickets)
			mux.Get("/events", s.myEvents)
			mux.Route("/profile", func(mux chi.Router) {
				mux.Get("/", s.profile)
				mux.Put("/", s.updateProfile)
			})
		})

	})

	s.mux = mux
}
