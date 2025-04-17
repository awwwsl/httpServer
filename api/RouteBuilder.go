package api

import (
	"httpServer/services"
	"net/http"
)

type RouteBuilder struct {
	Mux             *http.ServeMux
	ServiceProvider *services.ServiceProvider
}

func NewRouteBuilder(serviceProvider *services.ServiceProvider) *RouteBuilder {
	builder := &RouteBuilder{
		Mux:             http.NewServeMux(),
		ServiceProvider: serviceProvider,
	}
	return builder
}
