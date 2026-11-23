package marketdatapostgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/eternalpaininside/Algotrade.git/backend/internal/core/domain"
	coreerrors "github.com/eternalpaininside/Algotrade.git/backend/internal/core/errors"
	"github.com/jackc/pgx/v5"
)

func (r *CandleRepository) GetCandles(
	ctx context.Context,
	exchangeID int64,
	symbol, interval string,
	from, to time.Time,
) ([]domain.Candle, error) {
	seconds, err := secondsFromString(interval)
	if err != nil {
		return nil, err
	}
	duration := time.Duration(seconds) * time.Second

	query := `
	SELECT
		c.open_time,
		c.open_price,
		c.high_price,
		c.low_price,
		c.close_price,
		c.volume
	FROM algotrade.instruments i
	JOIN algotrade.candles c
	ON c.instrument_id = i.instrument_id
	WHERE i.exchange_id = $1
		AND i.symbol = $2
		AND c.interval_seconds = $3
		AND c.open_time >= $4
		AND c.open_time < $5
	ORDER BY c.open_time
`

	rows, err := r.pool.Query(
		ctx,
		query,
		exchangeID,
		duration,
		seconds,
		from,
		to,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf(
				"candles on exchange %d: %w",
				exchangeID,
				coreerrors.ErrNotFound,
			)
		}
		return nil, fmt.Errorf(
			"find candles: %w",
			err,
		)
	}

	candles := make([]domain.Candle, 0)
	for rows.Next() {
		candle := domain.Candle{
			ExchangeID: exchangeID,
			Symbol:     symbol,
			Interval:   interval,
		}
		if err = rows.Scan(
			&candle.OpenedAt,
			&candle.OpenPrice,
			&candle.HighPrice,
			&candle.LowPrice,
			&candle.ClosePrice,
			&candle.Volume,
		); err != nil {
			return nil, fmt.Errorf(
				"scan candle: %w",
				err,
			)
		}

		candle.OpenedAt = candle.OpenedAt.UTC()
		candle.ClosedAt = candle.OpenedAt.Add(
			duration - time.Millisecond,
		)
		candles = append(candles, candle)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"read candles: %w",
			err,
		)
	}
	return candles, nil
}
