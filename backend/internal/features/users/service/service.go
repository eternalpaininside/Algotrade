package usersservice

import (
	"context"

	"github.com/eternalpaininside/Algotrade.git/backend/internal/core/domain"
)

type UserService struct {
	repository UserRepository
}

type UserRepository interface {
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
	PatchUser(
		ctx context.Context,
		id int,
		user domain.User,
	) (domain.User, error)
}

func NewUserService(
	repository UserRepository,
) *UserService {
	return &UserService{
		repository: repository,
	}
}
