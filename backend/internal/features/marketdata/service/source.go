package marketservice

import (
	"context"
	"time"

	"github.com/eternalpaininside/Algotrade.git/backend/internal/core/domain"
)

type CandleSource interface {
	GetCandles(
		ctx context.Context,
		symbol string,
		interval string,
		from, to time.Time,
	) ([]domain.Candle, error)
}
