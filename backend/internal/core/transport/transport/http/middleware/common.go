package middleware

import (
	"context"
	"net/http"
	"strings"
	"uuid"

	core_logger "github.com/eternalpaininside/Algotrade.git/backend/internal/core/logger"
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

func Logging(logger *core_logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(RequestIDHeader)

			logger.With(
				zap.String("request_id", requestID),
				zap.String("url", r.URL.String()),
			)

			ctx := context.WithValue(
				r.Context(),
				"logger",
				logger,
			)

			next.ServeHTTP(rw, r.WithContext(ctx))
		})
	}
}

func Panic() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			
		})
	}
}
