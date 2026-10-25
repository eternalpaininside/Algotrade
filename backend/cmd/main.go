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

	logger, err := corelogger.NewLogger(corelogger.NewConfigMust())
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
	)

	usersRep := userspostgres.NewUserRepository(pool)
	usersService := usersservice.NewUserService(usersRep)
	usersTransport := usershttp.NewUserHTTPHandler(usersService)

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
	apiVersionRouter.RegisterRoutes(usersTransport.Routes()...)
	server.RegisterAPIRoutes(apiVersionRouter)

	if err = server.Run(ctx); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server run",
				zap.Error(err),
			)
		}
	}
}
