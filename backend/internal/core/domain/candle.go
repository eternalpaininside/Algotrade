package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type Candle struct {
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

func (c *Candle) Validate() error {
	if c.ExchangeID <= 0 {
		return fmt.Errorf(
			"negative exchange id: %d",
			c.ExchangeID,
		)
	}

	if strings.TrimSpace(c.Symbol) == "" {
		return errors.New("empty candle symbol")
	}

	if strings.TrimSpace(c.Interval) == "" {
		return errors.New("empty candle interval")
	}

	if !c.OpenedAt.Before(c.ClosedAt) {
		return fmt.Errorf(
			"invalid candle period: opened at %v and closed at %v",
			c.OpenedAt,
			c.ClosedAt,
		)
	}

	if c.OpenPrice < 0 ||
		c.HighPrice < 0 || c.LowPrice < 0 ||
		c.ClosePrice < 0 {
		return errors.New("negative price")
	}

	if c.Volume < 0 {
		return errors.New("negarive volume")
	}

	if c.LowPrice > c.OpenPrice ||
		c.LowPrice > c.ClosePrice ||
		c.LowPrice > c.HighPrice {
		return errors.New("low price more than other prices")
	}

	if c.HighPrice < c.OpenPrice ||
		c.HighPrice < c.ClosePrice ||
		c.HighPrice < c.LowPrice {
		return errors.New("high price lower than other prices")
	}
	return nil
}
