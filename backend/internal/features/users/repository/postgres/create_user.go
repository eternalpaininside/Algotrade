package userpostgres

import (
	"context"
	"fmt"

	"github.com/eternalpaininside/Algotrade.git/backend/internal/core/domain"
)

func (u *UserRepository) CreateUser(
	ctx context.Context,
	user domain.User,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(
		ctx,
		u.pool.OpTimeout(),
	)
	defer cancel()

	query := `
INSERT INTO algotrade.users(user_name, email, password, created_at)
VALUES ($1, $2, $3, $4)
RETURNING id, user_name, email, created_at;
`
	row := u.pool.QueryRow(
		ctx, query,
		user.Name, user.Email, user.Pass, user.CreatedAt,
	)

	var userModel UserModel
	if err := row.Scan(
		&userModel.ID,
		&userModel.Name,
		&userModel.Email,
		&userModel.Password,
		userModel.CreateAt,
	); err != nil {
		return domain.User{}, fmt.Errorf(
			"row scan error: %w",
			err,
		)
	}

	userDomain := domain.NewUser(
		userModel.ID,
		userModel.CreateAt,
		userModel.Name,
		userModel.Email,
		userModel.Password,
	)

	return userDomain, nil
}
