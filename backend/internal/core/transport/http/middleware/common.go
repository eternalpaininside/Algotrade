package coremiddleware

import (
	"context"
	"net/http"
	"strings"
	"time"
	"uuid"

	"github.com/eternalpaininside/Algotrade.git/backend/internal/core/logger"
	"github.com/eternalpaininside/Algotrade.git/backend/internal/core/transport/http/response"
	"go.uber.org/zap"
)

const RequestIDHeader = "X-Request-ID"

func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(RequestIDHeader)
			if strings.TrimSpace(requestID) == "" {
				requestID = uuid.New().String()
			}

			r.Header.Set(RequestIDHeader, requestID)
			rw.Header().Set(RequestIDHeader, requestID)

			next.ServeHTTP(rw, r)
		})
	}
}

func Logging(logger *corelogger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(RequestIDHeader)

			requestLogger := logger.With(
				zap.String("request_id", requestID),
				zap.String("url", r.URL.String()),
			)

			ctx := context.WithValue(
				r.Context(),
				"logger",
				requestLogger,
			)

			next.ServeHTTP(rw, r.WithContext(ctx))
		})
	}
}

func Panic() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			logger := corelogger.FromContext(ctx)

			handler := corehttpresponse.NewHTTPResponseHandler(
				logger,
				rw,
			)

			defer func() {
				if p := recover(); p != nil {
					handler.PanicResponse(
						p,
						"unexpected panic during HTTP handler",
					)
				}
			}()

			next.ServeHTTP(rw, r)
		})
	}
}

func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := corelogger.FromContext(ctx)
			start := time.Now()
			rw := corehttpresponse.NewResponseWriter(w)

			log.Debug(
				">>> incoming HTTP request",
				zap.String("http_method", r.Method),
				zap.Time("time", start.UTC()),
			)

			next.ServeHTTP(rw, r)

			log.Debug(
				"<<< done HTTP request",
				zap.Int("status_code", rw.StatusCodeOrPanic()),
				zap.Duration("latency", time.Since(start)),
			)
		})
	}
}
