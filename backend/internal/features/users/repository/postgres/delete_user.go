package userpostgres

import (
	"context"
	"fmt"

	coreerrors "github.com/eternalpaininside/Algotrade.git/backend/internal/core/errors"
)

func (u *UserRepository) DeleteUser(
	ctx context.Context,
	id int,
) error {
	ctx, cancel := context.WithTimeout(
		ctx,
		u.pool.OpTimeout(),
	)
	defer cancel()

	query := `
DELETE FROM algotrade.users
WHERE user_id=$1;
`
	tag, err := u.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf(
			"exec query: %w",
			err,
		)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf(
			"user with id='%d': %w",
			id,
			coreerrors.ErrNotFound,
		)
	}
	return nil
}
