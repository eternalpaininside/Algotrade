package usersservice

import (
	"context"
	"fmt"

	"github.com/eternalpaininside/Algotrade.git/backend/internal/core/domain"
)

func (u *UserService) PatchUser(
	ctx context.Context,
	id int,
	patch domain.UserPatch,
) (domain.User, error) {
	user, err := u.repository.GetUser(
		ctx,
		id,
	)
	if err != nil {
		return domain.User{}, fmt.Errorf(
			"failed to get user: %w",
			err,
		)
	}

	if err = user.ApplyPatch(patch); err != nil {
		return domain.User{}, fmt.Errorf(
			"failed to apply patch: %w",
			err,
		)
	}


	patched, err := u.repository.PatchUser(
		ctx,
		id,
		user,
	)
	if err != nil {
		return domain.User{}, fmt.Errorf(
			"failed to patch user: %w",
			err,
		)
	}

	return patched, nil
}
