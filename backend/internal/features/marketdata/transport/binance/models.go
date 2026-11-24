package markettransportbinance

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/eternalpaininside/Algotrade.git/backend/internal/core/domain"
)

type klineRow []json.RawMessage

func (r klineRow) toDomain(
	exchangeID int64,
	symbol, interval string,
) (domain.Candle, error) {
	const requiredFields = 7

	if len(r) < requiredFields {
		return domain.Candle{}, fmt.Errorf(
			"unexpected Binance kline field count: %d",
			len(r),
		)
	}

	openedAt, err := parseUnixMilli(
		r[0],
		"open time",
	)
	if err != nil {
		return domain.Candle{}, err
	}

	openPrice, err := parsePrice(
		r[1],
		"open price",
	)
	if err != nil {
		return domain.Candle{}, err
	}

	highPrice, err := parsePrice(
		r[2],
		"high price",
	)
	if err != nil {
		return domain.Candle{}, err
	}

	lowPrice, err := parsePrice(
		r[3],
		"low price",
	)
	if err != nil {
		return domain.Candle{}, err
	}

	closePrice, err := parsePrice(
		r[4],
		"close price",
	)
	if err != nil {
		return domain.Candle{}, err
	}

	volume, err := parsePrice(
		r[5],
		"volume",
	)
	if err != nil {
		return domain.Candle{}, err
	}

	closedAt, err := parseUnixMilli(
		r[6],
		"close time",
	)
	if err != nil {
		return domain.Candle{}, err
	}

	candle := domain.Candle{
		ExchangeID: exchangeID,
		Symbol: symbol,
		Interval: interval,

		OpenedAt: openedAt,
		ClosedAt: closedAt,

		OpenPrice: openPrice,
		HighPrice: highPrice,
		LowPrice: lowPrice,
		ClosePrice: closePrice,
		Volume: volume,
	}

	if err = candle.Validate(); err != nil {
		return domain.Candle{}, fmt.Errorf(
			"validate Binance candle: %w",
			err,
		)
	}

	return candle, nil
}

func parseUnixMilli(
	data json.RawMessage,
	field string,
) (time.Time, error) {
	var milliseconds int64

	if err := json.Unmarshal(
		data,
		&milliseconds,
	); err != nil {
		return time.Time{}, fmt.Errorf(
			"unmarshal Binance %s: %w",
			field,
			err,
		)
	}

	return time.UnixMilli(milliseconds).UTC(), nil
}

func parsePrice(
	data json.RawMessage,
	field string,
) (float64, error) {
	var value string

	if err := json.Unmarshal(
		data,
		&value,
	); err != nil {
		return 0, fmt.Errorf(
			"unmarshal Binance %s: %w",
			field,
			err,
		)
	}

	price, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, fmt.Errorf(
			"parse Binance %s: %w",
			field,
			err,
		)
	}

	return price, nil
}
