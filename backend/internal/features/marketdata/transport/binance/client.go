package markettransportbinance

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	defaultBaseURL = "https://api.binance.com"
	defaultTimeout = 10 * time.Second
)

type Config struct {
	ExchangeID int64
	BaseURL    string
	Timeout    time.Duration
}

type Client struct {
	httpClient *http.Client
	baseURL    *url.URL
	exchangeID int64
}

func NewClient(
	config Config,
) (*Client, error) {
	if config.ExchangeID <= 0 {
		return nil, fmt.Errorf(
			"invalid exchange id: %d",
			config.ExchangeID,
		)
	}

	baseURL := strings.TrimSpace(config.BaseURL)
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	parsedBaseURL, err := url.ParseRequestURI(baseURL)
	if err != nil {
		return nil, fmt.Errorf(
			"parse Binance base URL: %w",
			err,
		)
	}

	if parsedBaseURL.Scheme == "" ||
		parsedBaseURL.Host == "" {
		return nil, fmt.Errorf(
			"invalid Binance base URL: %q",
			baseURL,
		)
	}

	timeout := config.Timeout
	if timeout <= 0 {
		timeout = defaultTimeout
	}

	return &Client{
		httpClient: &http.Client{
			Timeout: timeout,
		},
		baseURL:    parsedBaseURL,
		exchangeID: config.ExchangeID,
	}, nil
}
