package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	corelogger "github.com/eternalpaininside/Algotrade.git/backend/internal/core/logger"
	coremiddleware "github.com/eternalpaininside/Algotrade.git/backend/internal/core/transport/http/middleware"
	corehttpserver "github.com/eternalpaininside/Algotrade.git/backend/internal/core/transport/http/server"
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

	logger.Debug("starting application")

	usersTransport := usershttp.NewUserHTTPHandler(nil)
	routes := usersTransport.Routes()

	apiVersionRouter := corehttpserver.NewApiRouter(
		corehttpserver.ApiVersion1,
	)
	apiVersionRouter.RegisterRoutes(routes...)

	server := corehttpserver.NewHTTPServer(
		corehttpserver.NewConfigMust(),
		logger,
		coremiddleware.RequestID(),
		coremiddleware.Logging(logger),
		coremiddleware.Panic(),
		coremiddleware.Trace(),
	)

	server.RegisterAPIRoutes(apiVersionRouter)

	if err = server.Run(ctx); err != nil {
		logger.Error("server run",
			zap.Error(err),
		)
	}
}
