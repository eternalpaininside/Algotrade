package markethttp

import (
	"context"
	"net/http"
	"time"

	corehttpserver "github.com/eternalpaininside/Algotrade.git/backend/internal/core/transport/http/server"
)

type MarketHTTPHandler struct {
	marketService MarketDataService
}

type MarketDataService interface {
	CollectCandles(
		ctx context.Context,
		symbol, interval string,
		from, to time.Time,
	) (int, error)
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
			Method:  http.MethodPost,
			Path:    "/market-data/candles/collect",
			Handler: m.CollectCandles,
		},
	}
}
