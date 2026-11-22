package markethttp

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	coreerrors "github.com/eternalpaininside/Algotrade.git/backend/internal/core/errors"
	corelogger "github.com/eternalpaininside/Algotrade.git/backend/internal/core/logger"
	corehttpresponse "github.com/eternalpaininside/Algotrade.git/backend/internal/core/transport/http/response"
	"go.uber.org/zap"
)

const (
	SymbolKey   = "symbol"
	IntervalKey = "interval"
	FromKey     = "from"
	ToKey       = "to"
)

type CollectCandlesResponse struct {
	CollectedCount int `json:"collected_count"`
}

func (h *MarketHTTPHandler) GetCandles(
	rw http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()
	log := corelogger.FromContext(ctx)
	responseHandler := corehttpresponse.NewHTTPResponseHandler(
		log,
		rw,
	)

	param, err := getCandleQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get candle query params",
		)
		return
	}

	log.Debug(
		"start collect market candles handler",
		zap.String("symbol", param.symbol),
		zap.String("interval", param.interval),
		zap.Time("from", param.from),
		zap.Time("to", param.to),
	)

	n, err := h.marketService.CollectCandles(
		ctx,
		param.symbol,
		param.interval,
		param.from,
		param.to,
	)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to collect market candles",
		)
		return
	}

	responseHandler.JSONResponse(
		CollectCandlesResponse{
			CollectedCount: n,
		},
		http.StatusOK,
	)
}

type candleQueryParams struct {
	symbol, interval string
	from, to         time.Time
}

func (c candleQueryParams) Validate() error {
	if c.symbol == "" {
		return errors.New("empty candle symbol")
	}

	if c.interval == "" {
		return errors.New("empty query interval")
	}

	if !c.from.Before(c.to) {
		return fmt.Errorf(
			"invalid candle period from %v to %v",
			c.from,
			c.to,
		)
	}

	return nil
}

func getCandleQueryParams(
	r *http.Request,
) (candleQueryParams, error) {
	var defaultParams candleQueryParams

	query := r.URL.Query()

	symbol := strings.TrimSpace(
		query.Get(SymbolKey),
	)
	interval := strings.TrimSpace(
		query.Get(IntervalKey),
	)

	from, err := parseQueryTime(
		query.Get(FromKey),
		FromKey,
	)
	if err != nil {
		return defaultParams, err
	}

	to, err := parseQueryTime(
		query.Get(ToKey),
		ToKey,
	)
	if err != nil {
		return defaultParams, err
	}

	param := candleQueryParams{
		symbol:   symbol,
		interval: interval,
		from:     from,
		to:       to,
	}

	if err = param.Validate(); err != nil {
		return defaultParams, fmt.Errorf(
			"candle query params validate: %v: %w",
			err,
			coreerrors.ErrInvalidArgument,
		)
	}
	return param, nil
}

func parseQueryTime(
	value, key string,
) (time.Time, error) {
	value = strings.TrimSpace(value)

	if value == "" {
		return time.Time{}, fmt.Errorf(
			"empty query param %q: %w",
			key,
			coreerrors.ErrInvalidArgument,
		)
	}

	parsedTime, err := time.Parse(
		time.RFC3339,
		value,
	)
	if err != nil {
		return time.Time{}, fmt.Errorf(
			"parse query patam %q=%q: %v: %w",
			key,
			value,
			err,
			coreerrors.ErrInvalidArgument,
		)
	}

	return parsedTime.UTC(), nil
}

