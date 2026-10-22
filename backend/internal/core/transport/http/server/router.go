package corehttpserver

import (
	"github.com/go-chi/chi/v5"
)

type ApiVersion string

const (
	ApiVersion1 = ApiVersion("v1")
	ApiVersion2 = ApiVersion("v2")
)

type APIVersionRouter struct {
	*chi.Mux
	apiVersion ApiVersion
}

func NewApiRouter(
	version ApiVersion,
) *APIVersionRouter {
	return &APIVersionRouter{
		Mux:        chi.NewMux(),
		apiVersion: version,
	}
}

func (r *APIVersionRouter) RegisterRoutes(
	routes ...Route,
) {
	for _, route := range routes {
		r.MethodFunc(
			route.Method,
			route.Path,
			route.Handler,
		)
	}
}
