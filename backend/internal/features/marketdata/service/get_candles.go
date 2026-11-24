package marketservice

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/eternalpaininside/Algotrade.git/backend/internal/core/domain"
	coreerrors "github.com/eternalpaininside/Algotrade.git/backend/internal/core/errors"
)

func (s *MarketDataService) GetCandles(
	ctx context.Context,
	exchangeID int64,
	symbol, interval string,
	from, to time.Time,
) ([]domain.Candle, error) {
	if exchangeID <= 0 {
		return nil, fmt.Errorf(
			"exchange id must be positive, got %d: %w",
			exchangeID,
			coreerrors.ErrInvalidArgument,
		)
	}

	symbol = strings.ToUpper(
		strings.TrimSpace(symbol),
	)
	if symbol == "" {
		return nil, fmt.Errorf(
			"empty candle symbol: %w",
			coreerrors.ErrInvalidArgument,
		)
	}

	interval = strings.TrimSpace(interval)
	if interval == "" {
		return nil, fmt.Errorf(
			"empty candle interval: %w",
			coreerrors.ErrInvalidArgument,
		)
	}

	if !from.Before(to) {
		return nil, fmt.Errorf(
			"invalid candle period from %v to %v: %w",
			from,
			to,
			coreerrors.ErrInvalidArgument,
		)
	}

	candles, err := s.repository.GetCandles(
		ctx,
		exchangeID,
		symbol, interval,
		from, to,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"get candles from repository: %w",
			err,
		)
	}

	return candles, nil
}
