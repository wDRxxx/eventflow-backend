package httpServer

import (
	"github.com/wDRxxx/eventflow-backend/internal/api"
	"github.com/wDRxxx/eventflow-backend/internal/service"
	"net/http"
)

type server struct {
	mux        http.Handler
	apiService service.ApiService
}

func NewHTTPServer(apiService service.ApiService) api.HTTPServer {
	s := &server{
		apiService: apiService,
	}
	s.setRoutes()

	return s
}

func (s *server) Handler() http.Handler {
	return s.mux
}
