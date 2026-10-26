package usersservice

import (
	"context"
	"fmt"

	"github.com/eternalpaininside/Algotrade.git/backend/internal/core/domain"
)

func (u *UserService) CreateUser(
	ctx context.Context,
	password string,
	user domain.User,
) (domain.User, error) {
	salt, hash, err := domain.PasswordHashing(password)
	if err != nil {
		return domain.User{}, fmt.Errorf(
			"hash password: %w",
			err,
		)
	}
	user.PassHash = domain.ConcatHash(salt, hash)

	if err = user.Validate(); err != nil {
		return domain.User{}, fmt.Errorf(
			"validate user domain: %w",
			err,
		)
	}

	user, err = u.repository.CreateUser(
		ctx,
		user,
	)
	if err != nil {
		return domain.User{}, fmt.Errorf(
			"repository create user: %w",
			err,
		)
	}

	return user, nil
}
