package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	corelogger "github.com/eternalpaininside/Algotrade.git/backend/internal/core/logger"
	postgrepool "github.com/eternalpaininside/Algotrade.git/backend/internal/core/repository/postgres/pool"
	coremiddleware "github.com/eternalpaininside/Algotrade.git/backend/internal/core/transport/http/middleware"
	corehttpserver "github.com/eternalpaininside/Algotrade.git/backend/internal/core/transport/http/server"
	marketdatapostgres "github.com/eternalpaininside/Algotrade.git/backend/internal/features/marketdata/repository/postgre"
	marketservice "github.com/eternalpaininside/Algotrade.git/backend/internal/features/marketdata/service"
	marketbinance "github.com/eternalpaininside/Algotrade.git/backend/internal/features/marketdata/transport/binance"
	markethttp "github.com/eternalpaininside/Algotrade.git/backend/internal/features/marketdata/transport/http"
	userspostgres "github.com/eternalpaininside/Algotrade.git/backend/internal/features/users/repository/postgres"
	usersservice "github.com/eternalpaininside/Algotrade.git/backend/internal/features/users/service"
	usershttp "github.com/eternalpaininside/Algotrade.git/backend/internal/features/users/transport/http"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	logger, err := corelogger.NewLogger(
		corelogger.NewConfigMust(),
	)
	if err != nil {
		log.Fatalf(
			"failed to init application logger: %v\n",
			err,
		)
	}
	defer logger.Close()

	logger.Debug("initializing pg connection pool")

	pool, err := postgrepool.NewConnectionPool(
		ctx,
		postgrepool.NewConfigMust(),
	)
	if err != nil {
		logger.Fatal(
			"failed to init postgres connection init",
			zap.Error(err),
		)
	}
	defer pool.Close()

	logger.Debug(
		"initializing feature",
		zap.String("feature", "users"),
		zap.String("feature", "marketdata"),
	)

	usersRep := userspostgres.NewUserRepository(pool)
	usersService := usersservice.NewUserService(usersRep)
	usersTransport := usershttp.NewUserHTTPHandler(usersService)

	binanceConfig := NewBinanceConfigMust(
		ctx,
		pool,
		logger,
	)
	binanceClient, err := marketbinance.NewClient(
		binanceConfig,
	)
	if err != nil {
		log.Fatal(
			"failed to create Binance client",
			zap.Error(err),
		)
	}

	candleRepository := marketdatapostgres.NewCandleRepository(pool)
	marketDataService := marketservice.NewMarketDataService(
		binanceClient,
		candleRepository,
	)
	marketTransport := markethttp.NewMarketHTTPHandler(
		marketDataService,
		binanceConfig.ExchangeID,
	)

	logger.Debug("initializing HTTP server")

	server := corehttpserver.NewHTTPServer(
		corehttpserver.NewConfigMust(),
		logger,
		coremiddleware.RequestID(),
		coremiddleware.Logging(logger),
		coremiddleware.Panic(),
		coremiddleware.Trace(),
	)

	apiVersionRouter := corehttpserver.NewApiRouter(
		corehttpserver.ApiVersion1,
	)
	apiVersionRouter.RegisterRoutes(
		usersTransport.Routes()...,
	)
	apiVersionRouter.RegisterRoutes(
		marketTransport.Routes()...,
	)

	server.RegisterAPIRoutes(
		apiVersionRouter,
	)

	if err = server.Run(ctx); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			logger.Error(
				"server run",
				zap.Error(err),
			)
		}
	}
}

func NewBinanceConfigMust(
	ctx context.Context,
	pool *postgrepool.ConnectionPool,
	logger *corelogger.Logger,
) marketbinance.Config {
	query := `
	SELECT exchange_id
	FROM algotrade.exchanges
	WHERE code = $1;
	`

	var exchangeID int64
	lookupCtx, cancel := context.WithTimeout(
		ctx,
		pool.OpTimeout(),
	)
	if err := pool.QueryRow(
		lookupCtx,
		query,
		"binance",
	).Scan(&exchangeID); err != nil {
		logger.Fatal(
			"failed to find Binance exchange",
			zap.Error(err),
		)
	}
	cancel()

	return marketbinance.Config{
		ExchangeID: exchangeID,
	}
}
