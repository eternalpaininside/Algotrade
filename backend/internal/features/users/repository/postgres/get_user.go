package userpostgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/eternalpaininside/Algotrade.git/backend/internal/core/domain"
	coreerrors "github.com/eternalpaininside/Algotrade.git/backend/internal/core/errors"
	"github.com/jackc/pgx/v5"
)

func (u *UserRepository) GetUser(
	ctx context.Context,
	id int,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(
		ctx,
		u.pool.OpTimeout(),
	)
	defer cancel()

	query := `
SELECT user_id, user_name, email, password, created_at
FROM algotrade.users
WHERE user_id=$1;
`

	row := u.pool.QueryRow(ctx, query, id)

	var userModel UserModel
	err := row.Scan(
		&userModel.ID,
		&userModel.Name,
		&userModel.Email,
		&userModel.Password,
		&userModel.CreateAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, fmt.Errorf(
				"user with id='%d': %w",
				id,
				coreerrors.ErrNotFound,
			)
		}
		return domain.User{}, fmt.Errorf(
			"query row scan error: %w",
			err,
		)
	}

	return domain.NewUser(
		userModel.ID,
		userModel.CreateAt,
		userModel.Name,
		userModel.Email,
		userModel.Password,
	), nil
}
