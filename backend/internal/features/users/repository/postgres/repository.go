package userpostgres

import postgrepool "github.com/eternalpaininside/Algotrade.git/backend/internal/core/repository/postgres/pool"

type UserRepository struct {
	pool postgrepool.Pool
}

func NewUserRepository(
	pool postgrepool.Pool,
) *UserRepository {
	return &UserRepository{
		pool: pool,
	}
}
