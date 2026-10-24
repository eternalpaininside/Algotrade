package corehttpserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	corelogger "github.com/eternalpaininside/Algotrade.git/backend/internal/core/logger"
	"github.com/eternalpaininside/Algotrade.git/backend/internal/core/transport/http/middleware"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type HTTPServer struct {
	mux    *chi.Mux
	config Config
	log    *corelogger.Logger

	middlewares []coremiddleware.Middleware
}

func NewHTTPServer(
	config Config,
	logger *corelogger.Logger,
	middlewares ...coremiddleware.Middleware,
) *HTTPServer {
	return &HTTPServer{
		mux:         chi.NewMux(),
		config:      config,
		log:         logger,
		middlewares: middlewares,
	}
}

func (h *HTTPServer) RegisterAPIRoutes(
	routers ...*APIVersionRouter,
) {
	for _, router := range routers {
		prefix := fmt.Sprintf(
			"/api/%s",
			router.apiVersion,
		)

		h.mux.Mount(
			prefix,
			router,
		)
	}
}

func (h *HTTPServer) Run(
	ctx context.Context,
) error {
	mux := coremiddleware.ChainMiddleware(
		h.mux,
		h.middlewares...,
	)

	server := &http.Server{
		Addr:    h.config.Address,
		Handler: mux,
	}

	chErr := make(chan error, 1)

	go func() {
		defer close(chErr)

		h.log.Warn(
			"start HTTP server",
			zap.String("addr", h.config.Address),
		)

		err := server.ListenAndServe()

		if !errors.Is(err, http.ErrServerClosed) {
			chErr <- err
		}
	}()

	select {
	case err, ok := <-chErr:
		if !ok {
			break
		}

		if err != nil {
			return fmt.Errorf(
				"listen and serve HTTP: %w",
				err,
			)
		}
	case <-ctx.Done():
		h.log.Warn("shutdown HTTP server")

		shutdownContext, cancel := context.WithTimeout(
			context.Background(),
			h.config.ShutdownTimeout,
		)
		defer cancel()

		if err := server.Shutdown(shutdownContext); err != nil {
			_ = server.Close()

			return fmt.Errorf(
				"shutdown HTTP server: %w",
				err,
			)
		}

		h.log.Warn("HTTP server stopped")
	}

	return nil
}
