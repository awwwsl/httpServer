package main

import (
	"httpServer/logging"
	"net/http"
)

type loggingServeMux struct {
	serverMux *http.ServeMux
	log       logging.ILogger
}

func newLoggingServeMux(log logging.ILogger, mux *http.ServeMux) *loggingServeMux {
	return &loggingServeMux{
		serverMux: mux,
		log:       log,
	}
}

func (l *loggingServeMux) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// Log and redirect to inner ServeMux
	l.log.Verbose("Received request %s from %s to url %s", request.Method, request.Host, request.RequestURI)
	l.serverMux.ServeHTTP(writer, request)
}
