package usersservice

import (
	"context"
	"fmt"

	"github.com/eternalpaininside/Algotrade.git/backend/internal/core/domain"
)

func (u *UserService) GetUser(
	ctx context.Context,
	id int,
) (domain.User, error) {
	user, err := u.repository.GetUser(ctx, id)
	if err != nil {
		return domain.User{}, fmt.Errorf(
			"repository get user with id='%d': %w",
			id,
			err,
		)
	}

	return user, nil
}
