package coremiddleware

import "slices"

import "net/http"

type Middleware func(handler http.Handler) http.Handler

func ChainMiddleware(
	handler http.Handler,
	middlewares ...Middleware,
) http.Handler {
	if len(middlewares) == 0 {
		return handler
	}

	for _, middleware := range slices.Backward(
		middlewares,
	) {
		handler = middleware(handler)
	}

	return handler
}
