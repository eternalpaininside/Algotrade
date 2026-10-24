package usersservice

import (
	"context"
	"fmt"

	"github.com/eternalpaininside/Algotrade.git/backend/internal/core/domain"
)

func (u *UserService) CreateUser(
	ctx context.Context,
	user domain.User,
) (domain.User, error) {
	if err := user.Validate(); err != nil {
		return domain.User{}, fmt.Errorf(
			"validate user domain: %w",
			err,
		)
	}

	user, err := u.repository.CreateUser(ctx, user)
	if err != nil {
		return domain.User{}, fmt.Errorf(
			"repository create user: %w",
			err,
		)
	}

	return user, nil
}
