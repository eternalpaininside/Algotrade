package corehttpresponse

import "net/http"

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

func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.ResponseWriter.WriteHeader(statusCode)
	rw.statusCode = statusCode
}

func (rw *ResponseWriter) StatusCodeOrPanic() int {
	return rw.statusCode
}
