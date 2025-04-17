package api

import "github.com/swaggest/openapi-go/openapi3"

type OpenApiBuilder struct {
	OpenApiReflector *openapi3.Reflector
}

func NewOpenApiBuilder() *OpenApiBuilder {
	return &OpenApiBuilder{
		OpenApiReflector: openapi3.NewReflector(),
	}
}
