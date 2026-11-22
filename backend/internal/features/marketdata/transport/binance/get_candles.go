package markettransportbinance

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/eternalpaininside/Algotrade.git/backend/internal/core/domain"
)

const maxKlinesPerReuest = 1000

func (c *Client) GetCandles(
	ctx context.Context,
	symbol, interval string,
	from, to time.Time,
) ([]domain.Candle, error) {
	symbol = strings.ToUpper(
		strings.TrimSpace(symbol),
	)
	interval = strings.TrimSpace(interval)

	if symbol == "" {
		return nil, errors.New("empty Binance symbol")
	}
	if interval == "" {
		return nil, errors.New("empty Binance internal")
	}

	if to.Before(from) {
		return nil, fmt.Errorf(
			"invalid candle period: from %v and to %v",
			from,
			to,
		)
	}

	from, to = from.UTC(), to.UTC()

	var candles []domain.Candle
	cursor := from

	for cursor.Before(to) {
		rows, err := c.getKlines(
			ctx,
			symbol,
			interval,
			cursor,
			to,
		)
		if err != nil {
			return nil, err
		}

		if len(rows) == 0 {
			break
		}

		var lastOpenedAt time.Time

		for _, row := range rows {
			candle, err := row.toDomain(
				c.exchangeID,
				symbol,
				interval,
			)
			if err != nil {
				return nil, err
			}

			if candle.OpenedAt.Before(from) ||
				!candle.OpenedAt.Before(to) {
				continue
			}

			candles = append(
				candles,
				candle,
			)
			lastOpenedAt = candle.OpenedAt
		}

		if len(rows) < maxKlinesPerReuest {
			break
		}

		nextCursor := lastOpenedAt.Add(time.Millisecond)
		if cursor.After(nextCursor) {
			return nil, errors.New("Binance kline pagination did not advance")
		}

		cursor = nextCursor
	}

	return candles, nil
}

func (c *Client) getKlines(
	ctx context.Context,
	symbol, interval string,
	from, to time.Time,
) ([]klineRow, error) {
	endpoint := c.baseURL.JoinPath(
		"api",
		"v3",
		"klines",
	)

	var query url.Values
	query.Set("symbol", symbol)
	query.Set("interval", interval)
	query.Set(
		"startTime",
		strconv.FormatInt(
			from.UnixMilli(),
			10,
		),
	)
	query.Set(
		"endTime",
		strconv.FormatInt(
			to.Add(-time.Millisecond).UnixMilli(),
			10,
		),
	)
	query.Set(
		"limit",
		strconv.Itoa(maxKlinesPerReuest),
	)

	endpoint.RawQuery = query.Encode()

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		endpoint.String(),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"create Binance klines request: %w",
			err,
		)
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf(
			"do Binance klines request: %w",
			err,
		)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, err := io.ReadAll(
			io.LimitReader(request.Body, 4*1024),
		)
		if err != nil {
			return nil, fmt.Errorf(
				"read Binance error response: %w",
				err,
			)
		}

		return nil, fmt.Errorf(
			"Binance klines response status=%d, body=%q",
			response.StatusCode,
			string(body),
		)
	}

	var rows []klineRow

	if err = json.NewDecoder(
		request.Body,
	).Decode(&rows); err != nil {
		return nil, fmt.Errorf(
			"decode Binance klines response: %w",
			err,
		)
	}

	return rows, nil
}
