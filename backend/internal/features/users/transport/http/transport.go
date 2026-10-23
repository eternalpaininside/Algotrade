package usershttp

import (
	"context"
	"net/http"

	"github.com/eternalpaininside/Algotrade.git/backend/internal/core/domain"
	corehttpserver "github.com/eternalpaininside/Algotrade.git/backend/internal/core/transport/http/server"
)

type UserHTTPHandler struct {
	userService UserService
}

type UserService interface {
	CreateUser(
		ctx context.Context,
		user domain.User,
	) (domain.User, error)
}

func NewUserHTTPHandler(
	service UserService,
) *UserHTTPHandler {
	return &UserHTTPHandler{
		userService: service,
	}
}

func (u *UserHTTPHandler) Routes() []corehttpserver.Route {
	return []corehttpserver.Route{
		{
			Method:  http.MethodPost,
			Path:    "/users",
			Handler: u.CreateUser,
		},
	}
}
