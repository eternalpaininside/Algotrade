package markethttp

import (
	"context"
	"net/http"
	"time"

	"github.com/eternalpaininside/Algotrade.git/backend/internal/core/domain"
	corehttpserver "github.com/eternalpaininside/Algotrade.git/backend/internal/core/transport/http/server"
)

type MarketHTTPHandler struct{
	marketService MarketDataService
}

type MarketDataService interface {
	GetCandles(
		ctx context.Context,
		symbol, interval string,
		from, to time.Time,
	) ([]domain.Candle, error)
}

func NewMarketHTTPHandler(
	service MarketDataService,
) *MarketHTTPHandler {
	return &MarketHTTPHandler{
		marketService: service,
	}
}

func (m *MarketHTTPHandler) Routes() []corehttpserver.Route {
	return []corehttpserver.Route{
		{
			Method: http.MethodGet,
			Path: "/market-data/candles",
			Handler: m.GetCandles,
		},
	}
}
