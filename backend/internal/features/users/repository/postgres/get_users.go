package userpostgres

import (
	"context"
	"fmt"

	"github.com/eternalpaininside/Algotrade.git/backend/internal/core/domain"
)

func (r *UserRepository) GetUsers(
	ctx context.Context,
	limit, offset *int,
) ([]domain.User, error) {
	ctx, cancel := context.WithTimeout(
		ctx, r.pool.OpTimeout(),
	)
	defer cancel()

	query := `
SELECT id, name, email, pass, created_at
FROM algotrade.users
ORDER BY id ASC
LIMIT $1
OFFSET $2;`

	rows, err := r.pool.Query(
		ctx,
		query,
		limit,
		offset,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"select users: %w",
			err,
		)
	}
	defer rows.Close()

	userModels := make([]UserModel, 0)
	for rows.Next() {
		var userModel UserModel

		err = rows.Scan(
			&userModel.ID,
			&userModel.Name,
			&userModel.Email,
			&userModel.Password,
			&userModel.CreateAt,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"scan users: %w",
				err,
			)
		}

		userModels = append(userModels, userModel)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"next rows: %w",
			err,
		)
	}

	userDomains := userDomainsFromModels(userModels)
	return userDomains, nil
}
