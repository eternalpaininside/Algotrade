package postgrepool

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool interface {
	Query(
		ctx context.Context,
		sql string,
		args ...any,
	) (pgx.Rows, error)
	QueryRow(
		ctx context.Context,
		sql string,
		args ...any,
	) pgx.Row
	Exec(
		ctx context.Context,
		sql string,
		arguments ...any,
	) (pgconn.CommandTag, error)
	Close()

	OpTimeout() time.Duration
}

type ConnectionPool struct {
	*pgxpool.Pool
	opTimeout time.Duration
}

func NewConnectionPool(
	ctx context.Context,
	config PoolConfig,
) (*ConnectionPool, error) {
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.PGUser,
		config.PGPassword,
		config.Host,
		config.Port,
		config.PGDatabase,
	)

	pgxConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf(
			"parse pgxconfig: %w",
			err,
		)
	}

	pool, err := pgxpool.NewWithConfig(
		ctx,
		pgxConfig,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"create pgxpool: %w",
			err,
		)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf(
			"pgxpool ping: %w",
			err,
		)
	}
	return &ConnectionPool{
		Pool:      pool,
		opTimeout: config.Timeout,
	}, nil
}

func (p *ConnectionPool) OpTimeout() time.Duration {
	return p.opTimeout
}
