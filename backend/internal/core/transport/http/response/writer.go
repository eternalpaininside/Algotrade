package corehttpresponse

import (
	"fmt"
	"net/http"
)

const StatusCodeUninitialized = -1

type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(
	rw http.ResponseWriter,
) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: rw,
		statusCode:     StatusCodeUninitialized,
	}
}

func (rw *ResponseWriter) WriteHeader(
	statusCode int,
) {
	if rw.statusCode == StatusCodeUninitialized {
		rw.ResponseWriter.WriteHeader(
			statusCode,
		)
		rw.statusCode = statusCode
	}
}

func (rw *ResponseWriter) StatusCodeOrPanic() int {
	if rw.statusCode == StatusCodeUninitialized {
		panic(
			fmt.Errorf(
				"uninitialized status code: %d",
				rw.statusCode,
			),
		)
	}
	return rw.statusCode
}
