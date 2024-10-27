package server

import (
	"net/http"
)

type ResponseWriter interface {
	http.ResponseWriter

	// StatusCode returns the HTTP response status code of the current request.
	StatusCode() int
}

var _ ResponseWriter = (*response)(nil)

type response struct {
	http.ResponseWriter
	statusCode int
}

func (r *response) reset(writer http.ResponseWriter) {
	r.ResponseWriter = writer
	r.statusCode = http.StatusOK
}

func (r *response) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *response) StatusCode() int {
	return r.statusCode
}
