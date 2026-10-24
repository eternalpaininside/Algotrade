package usersservice

import (
	"context"
	"fmt"
)

func (u *UserService) DeleteUser(
	ctx context.Context,
	id int,
) error {
	if err := u.repository.DeleteUser(
		ctx, id,
	); err != nil {
		return fmt.Errorf(
			"repository delete user: %w",
			err,
		)
	}

	return nil
}
