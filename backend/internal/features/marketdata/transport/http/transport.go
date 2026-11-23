package markethttp

import (
	"context"
	"net/http"
	"time"

	"github.com/eternalpaininside/Algotrade.git/backend/internal/core/domain"
	corehttpserver "github.com/eternalpaininside/Algotrade.git/backend/internal/core/transport/http/server"
)

type MarketHTTPHandler struct {
	marketService MarketDataService
	exchangeID    int64
}

type MarketDataService interface {
	CollectCandles(
		ctx context.Context,
		symbol, interval string,
		from, to time.Time,
	) (int, error)
	GetCandles(
		ctx context.Context,
		exchangeID int64,
		symbol, interval string,
		from, to time.Time,
	) ([]domain.Candle, error)
}

func NewMarketHTTPHandler(
	service MarketDataService,
	exchangeID int64,
) *MarketHTTPHandler {
	return &MarketHTTPHandler{
		marketService: service,
		exchangeID:    exchangeID,
	}
}

func (m *MarketHTTPHandler) Routes() []corehttpserver.Route {
	return []corehttpserver.Route{
		{
			Method:  http.MethodPost,
			Path:    "/market-data/candles/collect",
			Handler: m.CollectCandles,
		},
		{
			Method:  http.MethodGet,
			Path:    "/market-data/candles",
			Handler: m.GetCandles,
		},
	}
}
