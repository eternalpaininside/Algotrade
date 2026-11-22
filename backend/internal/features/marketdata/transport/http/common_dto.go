package markethttp

import (
	"time"

	"github.com/eternalpaininside/Algotrade.git/backend/internal/core/domain"
)

type CandleDTO struct {
	ExchangeID int64
	Symbol     string
	Interval   string

	OpenedAt time.Time
	ClosedAt time.Time

	OpenPrice  float64
	HighPrice  float64
	LowPrice   float64
	ClosePrice float64
	Volume     float64
}

func candleDTOFromDomain(
	candle domain.Candle,
) CandleDTO {
	return CandleDTO{
		ExchangeID: candle.ExchangeID,
		Symbol: candle.Symbol,
		Interval: candle.Interval,
		
		OpenedAt: candle.OpenedAt,
		ClosedAt: candle.ClosedAt,

		OpenPrice: candle.OpenPrice,
		HighPrice: candle.HighPrice,
		LowPrice: candle.LowPrice,
		ClosePrice: candle.ClosePrice,
		Volume: candle.Volume,
	}
}

