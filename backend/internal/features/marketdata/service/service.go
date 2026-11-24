package marketservice

import (
	"context"
	"time"

	"github.com/eternalpaininside/Algotrade.git/backend/internal/core/domain"
)

type MarketDataService struct {
	source CandleSource
	repository CandleRepository
}

type CandleSource interface {
	GetCandles(
		ctx context.Context,
		symbol, interval string,
		from, to time.Time,
	) ([]domain.Candle, error)
}

type CandleRepository interface {
	SaveCandles(
		ctx context.Context,
		candels []domain.Candle,
	) (int, error) 
	GetCandles(
		ctx context.Context,
		exchangeID int64,
		symbol, interval string,
		from, to time.Time,
	) ([]domain.Candle, error)
}

func NewMarketDataService(
	source CandleSource,
	repository CandleRepository,
) *MarketDataService {
	return &MarketDataService{
		source: source,
		repository: repository,
	}
}
