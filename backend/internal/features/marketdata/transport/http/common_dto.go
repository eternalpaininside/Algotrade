package markethttp

import (
	"time"

	"github.com/eternalpaininside/Algotrade.git/backend/internal/core/domain"
)

type CandleDTO struct {
	ExchangeID int64  `json:"exchange_id"`
	Symbol     string `json:"symbol"`
	Interval   string `json:"inteval"`

	OpenedAt time.Time `json:"opened_at"`
	ClosedAt time.Time `json:"closed_at"`

	OpenPrice  float64 `json:"open_price"`
	HighPrice  float64 `json:"high_price"`
	LowPrice   float64 `json:"low_price"`
	ClosePrice float64 `json:"close_price"`
	Volume     float64 `json:"volume"`
}

func candleDTOFromDomain(
	candle domain.Candle,
) CandleDTO {
	return CandleDTO{
		ExchangeID: candle.ExchangeID,
		Symbol:     candle.Symbol,
		Interval:   candle.Interval,

		OpenedAt: candle.OpenedAt,
		ClosedAt: candle.ClosedAt,

		OpenPrice:  candle.OpenPrice,
		HighPrice:  candle.HighPrice,
		LowPrice:   candle.LowPrice,
		ClosePrice: candle.ClosePrice,
		Volume:     candle.Volume,
	}
}
