package httpServer

import (
	"net/http"

	"github.com/wDRxxx/eventflow-backend/internal/api"
	"github.com/wDRxxx/eventflow-backend/internal/config"
	"github.com/wDRxxx/eventflow-backend/internal/service"
)

type server struct {
	mux http.Handler

	apiService service.ApiService
	authConfig *config.AuthConfig
	httpConfig *config.HttpConfig
}

func NewHTTPServer(
	apiService service.ApiService,
	authConfig *config.AuthConfig,
	httpConfig *config.HttpConfig,
) api.HTTPServer {
	s := &server{
		apiService: apiService,
		authConfig: authConfig,
		httpConfig: httpConfig,
	}

	s.setRoutes()

	return s
}

func (s *server) Handler() http.Handler {
	return s.mux
}
