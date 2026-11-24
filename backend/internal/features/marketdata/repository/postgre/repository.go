package marketdatapostgres

import postgrepool "github.com/eternalpaininside/Algotrade.git/backend/internal/core/repository/postgres/pool"

type CandleRepository struct {
	pool postgrepool.Pool
}

func NewCandleRepository(
	pool postgrepool.Pool,
) *CandleRepository {
	return &CandleRepository{
		pool: pool,
	}
}
