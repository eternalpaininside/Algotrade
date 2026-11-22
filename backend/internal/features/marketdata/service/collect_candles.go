package marketservice

import (
	"context"
	"fmt"
	"strings"
	"time"

	coreerrors "github.com/eternalpaininside/Algotrade.git/backend/internal/core/errors"
)

func (m *MarketDataService) CollectCandles(
	ctx context.Context,
	symbol, interval string,
	from, to time.Time,
) (int, error) {
	symbol = strings.TrimSpace(symbol)
	if symbol == "" {
		return 0, fmt.Errorf(
			"empty candle symbol: %w",
			coreerrors.ErrInvalidArgument,
		)
	}

	interval = strings.TrimSpace(interval)
	if interval == "" {
		return 0, fmt.Errorf(
			"empty candle interval: %w",
			coreerrors.ErrInvalidArgument,
		)
	}

	if !from.Before(to) {
		return 0, fmt.Errorf(
			"invalid candle period from %v and to %v: %w",
			from,
			to,
			coreerrors.ErrInvalidArgument,
		)
	}

	candles, err := m.source.GetCandles(
		ctx,
		symbol,
		interval,
		from.UTC(),
		to.UTC(),
	)
	if err != nil {
		return 0, fmt.Errorf(
			"get candles from source: %w",
			err,
		)
	}

	savedCount, err := m.repository.SaveCandles(
		ctx,
		candles,
	)
	if err != nil {
		return 0, fmt.Errorf(
			"save candles in repository: %w",
			err,
		)
	}
	return savedCount, nil
}
