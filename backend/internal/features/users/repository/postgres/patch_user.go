package userpostgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/eternalpaininside/Algotrade.git/backend/internal/core/domain"
	coreerrors "github.com/eternalpaininside/Algotrade.git/backend/internal/core/errors"
	"github.com/jackc/pgx/v5"
)

func (u *UserRepository) PatchUser(
	ctx context.Context,
	id int,
	user domain.User,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(
		ctx,
		u.pool.OpTimeout(),
	)
	defer cancel()

	query := `
UPDATE algotrade.users
SET 
    user_name = $1,
    email = $2,
    password = $3
WHERE user_id = $4
RETURNING 
	user_id,
	user_name,
	email
	password,
	created_at;
`

	row := u.pool.QueryRow(
		ctx,
		query,
		user.Name,
		user.Email,
		id,
	)

	var userModel UserModel
	if err := row.Scan(
		&userModel.ID,
		&userModel.Name,
		&userModel.Email,
		&userModel.Password,
		&userModel.CreatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, fmt.Errorf(
				"user with id='%d' concurrently: %w",
				id,
				coreerrors.ErrConflict,
			)
		}
		return domain.User{}, fmt.Errorf(
			"failed to scan row: %w",
			err,
		)
	}

	userDomain := domain.NewUser(
		userModel.ID,
		userModel.CreatedAt,
		userModel.Name,
		userModel.Email,
		userModel.Password,
	)
	return userDomain, nil
}
