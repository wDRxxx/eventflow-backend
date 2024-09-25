package httpServer

import (
	"net/http"

	"github.com/wDRxxx/eventflow-backend/internal/api"
	"github.com/wDRxxx/eventflow-backend/internal/config"
	"github.com/wDRxxx/eventflow-backend/internal/service"
)

type server struct {
	mux     http.Handler
	origins []string

	apiService service.ApiService
	authConfig *config.AuthConfig
}

func NewHTTPServer(
	apiService service.ApiService,
	authConfig *config.AuthConfig,
	origins []string,
) api.HTTPServer {
	s := &server{
		apiService: apiService,
		authConfig: authConfig,
		origins:    origins,
	}

	s.setRoutes()

	return s
}

func (s *server) Handler() http.Handler {
	return s.mux
}
