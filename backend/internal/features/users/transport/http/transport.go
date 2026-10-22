package usershttp

import (
	"net/http"

	corehttpserver "github.com/eternalpaininside/Algotrade.git/backend/internal/core/transport/http/server"
)

type UserHTTPHandler struct {
	userService UserService
}

type UserService interface {
}

func NewUserHTTPHandler(
	service UserService,
) *UserHTTPHandler {
	return &UserHTTPHandler{
		userService: service,
	}
}

func (h *UserHTTPHandler) Routes() []corehttpserver.Route {
	return []corehttpserver.Route{
		{
			Method:  http.MethodPost,
			Path:    "/users",
			Handler: h.CreateUser,
		},
	}
}
