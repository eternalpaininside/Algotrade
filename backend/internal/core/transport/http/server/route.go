package corehttpserver

import "net/http"

type Route struct {
	Method  string
	Path    string
	Handler http.HandlerFunc
}

func NewRoute(
	method, path string,
	handler http.HandlerFunc,
) Route {
	return Route{
		Method:  method,
		Path:    path,
		Handler: handler,
	}
}
