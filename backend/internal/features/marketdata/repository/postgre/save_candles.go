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

func (r *CandleRepository) SaveCandles(
	ctx context.Context,
	candles []domain.Candle,
) (int, error) {
	if len(candles) == 0 {
		return 0, nil
	}

	first := candles[0]

	var intervalSeconds int64
	switch first.Interval {
	case "1h":
		intervalSeconds = 3600
	case "1d":
		intervalSeconds = 86400
	default:
		return 0, fmt.Errorf(
			"unsupported candle interval %q: %w",
			first.Interval,
			coreerrors.ErrInvalidArgument,
		)
	}

	openTimes := make([]time.Time, len(candles))
	openPrices := make([]float64, len(candles))
	highPrices := make([]float64, len(candles))
	lowPrices := make([]float64, len(candles))
	closePrices := make([]float64, len(candles))
	volumes := make([]float64, len(candles))

	for i, candle := range candles {
		if candle.ExchangeID != first.ExchangeID ||
			candle.Symbol != first.Symbol ||
			candle.Interval != first.Interval {
			return 0, fmt.Errorf(
				"candle at index %d belongs to a different series: %w",
				i,
				coreerrors.ErrInvalidArgument,
			)
		}

		openTimes[i] = candle.OpenedAt
		openPrices[i] = candle.OpenPrice
		highPrices[i] = candle.HighPrice
		lowPrices[i] = candle.LowPrice
		closePrices[i] = candle.ClosePrice
		volumes[i] = candle.Volume
	}

	ctx, cancel := context.WithTimeout(
		ctx,
		r.pool.OpTimeout(),
	)
	defer cancel()

	instrumentID, err := r.findInstrumentID(
		ctx,
		first.ExchangeID,
		first.Symbol,
	)
	if err != nil {
		return 0, err
	}

	affected, err := r.InsertCandles(
		ctx,
		instrumentID,
		intervalSeconds,
		openTimes,
		openPrices,
		highPrices,
		lowPrices,
		closePrices,
		volumes,
	)
	if err != nil {
		return 0, err
	}

	return affected, nil
}

func (r *CandleRepository) findInstrumentID(
	ctx context.Context,
	exchangeID int64,
	symbol string,
) (int64, error) {
	query := `
	SELECT instrument_id
	FROM algotrade.instruments
	WHERE exchange_id = $1 AND symbol = $2
`
	var instrumentID int64
	if err := r.pool.QueryRow(
		ctx,
		query,
		exchangeID,
		symbol,
	).Scan(&instrumentID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf(
				"instrument %q on exchange %d: %w",
				symbol,
				exchangeID,
				coreerrors.ErrNotFound,
			)
		}

		return 0, fmt.Errorf(
			"find candle instrument: %w",
			err,
		)
	}

	return instrumentID, nil
}

func (r *CandleRepository) InsertCandles(
	ctx context.Context,
	instrumentID, intervalSeconds int64,
	openTimes []time.Time,
	openPrices, highPrices, lowPrices []float64,
	closePrices, volumes []float64,
) (int, error) {
	query := `
	INSERT INTO algotrade.candles (
		instrument_id, interval_seconds, open_time,
		open_price, high_price, low_price, close_price,
		volume
	)
	SELECT $1::BIGINT, $2::INTEGER, batch.*
	FROM unnest(
		$3::TIMESTAMPTZ[],
		$4::DOUBLE PRECISION[],
		$5::DOUBLE PRECISION[],
		$6::DOUBLE PRECISION[],
		$7::DOUBLE PRECISION[],
		$8::DOUBLE PRECISION[]
	) AS bach(open_time, open_price, high_price, low_price,
	close_price, volume)
	ON CONFLICT (instrument_id, interval_seconds, open_time) DO NOTHING
`

	tag, err := r.pool.Exec(
		ctx,
		query,
		instrumentID,
		intervalSeconds,
		openTimes,
		openPrices,
		highPrices,
		lowPrices,
		closePrices,
		volumes,
	)
	if err != nil {
		return 0, fmt.Errorf(
			"insert candles: %w",
			err,
		)
	}

	return int(tag.RowsAffected()), nil
}
