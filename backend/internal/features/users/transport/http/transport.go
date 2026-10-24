package usershttp

import (
	"context"
	"net/http"

	"github.com/eternalpaininside/Algotrade.git/backend/internal/core/domain"
	corehttpserver "github.com/eternalpaininside/Algotrade.git/backend/internal/core/transport/http/server"
)

const (
	UserIDKey = "id"
)

type UserHTTPHandler struct {
	userService UserService
}

type UserService interface {
	CreateUser(
		ctx context.Context,
		user domain.User,
	) (domain.User, error)
	GetUsers(
		ctx context.Context,
		limit, offset *int,
	) ([]domain.User, error)
	GetUser(
		ctx context.Context,
		id int,
	) (domain.User, error)
	DeleteUser(
		ctx context.Context,
		id int,
	) error
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
		{
			Method:  http.MethodGet,
			Path:    "/users",
			Handler: u.GetUsers,
		},
		{
			Method:  http.MethodGet,
			Path:    "/users/{id}",
			Handler: u.GetUser,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/users/{id}",
			Handler: u.DeleteUser,
		},
		{
			Method:  http.MethodPatch,
			Path:    "/users/{id}",
			Handler: u.PatchUser,
		},
	}
}
