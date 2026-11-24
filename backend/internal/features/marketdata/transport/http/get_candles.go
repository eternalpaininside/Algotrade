package markethttp

import (
	"net/http"

	corelogger "github.com/eternalpaininside/Algotrade.git/backend/internal/core/logger"
	corehttpresponse "github.com/eternalpaininside/Algotrade.git/backend/internal/core/transport/http/response"
	"go.uber.org/zap"
)

func (h *MarketHTTPHandler) GetCandles(
	rw http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()
	log := corelogger.FromContext(ctx)
	response := corehttpresponse.NewHTTPResponseHandler(
		log,
		rw,
	)

	params, err := getCandleQueryParams(r)
	if err != nil {
		response.ErrorResponse(
			err,
			"failed to get candle query params",
		)
		return
	}

	log.Debug(
		"get stored candles",
		zap.Int64("exchange_id", h.exchangeID),
		zap.String("symbol", params.symbol),
		zap.String("interval", params.interval),
		zap.Time("from", params.from),
		zap.Time("to", params.to),
	)

	candels, err := h.marketService.GetCandles(
		ctx,
		h.exchangeID,
		params.symbol,
		params.interval,
		params.from,
		params.to,
	)
	if err != nil {
		response.ErrorResponse(
			err,
			"failed to get stored candles",
		)
		return
	}

	items := make([]CandleDTO, 0, len(candels))
	for _, candle := range candels {
		items = append(items, candleDTOFromDomain(candle))
	}

	response.JSONResponse(items, http.StatusOK)
}
