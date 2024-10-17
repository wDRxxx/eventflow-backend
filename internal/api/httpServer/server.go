package httpServer

import (
	"net/http"
	"strings"

	"github.com/wDRxxx/eventflow-backend/internal/api"
	"github.com/wDRxxx/eventflow-backend/internal/config"
	"github.com/wDRxxx/eventflow-backend/internal/models"
	"github.com/wDRxxx/eventflow-backend/internal/service"
	"github.com/wDRxxx/eventflow-backend/internal/utils"
)

type server struct {
	mux http.Handler

	authConfig *config.AuthConfig
	httpConfig *config.HttpConfig

	eventsService  service.EventsService
	ticketsService service.TicketsService
	usersService   service.UsersService
}

func NewHTTPServer(
	authConfig *config.AuthConfig,
	httpConfig *config.HttpConfig,

	eventsService service.EventsService,
	ticketsService service.TicketsService,
	usersService service.UsersService,
) api.HTTPServer {
	s := &server{
		authConfig:     authConfig,
		httpConfig:     httpConfig,
		eventsService:  eventsService,
		ticketsService: ticketsService,
		usersService:   usersService,
	}

	s.setRoutes()

	return s
}

func (s *server) Handler() http.Handler {
	return s.mux
}

func (s *server) getAndVerifyHeaderToken(r *http.Request) (string, *models.UserClaims, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", nil, errInvalidAuthHeader
	}

	exploded := strings.Split(authHeader, " ")
	if len(exploded) != 2 || exploded[0] != "Bearer" {
		return "", nil, errInvalidAuthHeader
	}

	token := exploded[1]
	claims, err := utils.VerifyToken(token, s.authConfig.AccessTokenSecret())
	if err != nil {
		return "", nil, errInvalidAuthHeader
	}

	return token, claims, nil
}
