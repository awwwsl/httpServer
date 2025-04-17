package api

import (
	"github.com/swaggest/openapi-go"
	"net/http"
)

func RouteOpenApiFile(path string, route *RouteBuilder, openapi *OpenApiBuilder) {
	route.Mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS") // TODO: Dev only
			return
		} else if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		marshalJSON, err := openapi.OpenApiReflector.Spec.MarshalJSON()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		//w.Header().Set("Access-Control-Allow-Origin", "*")                   // TODO: Dev only
		//w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS") // TODO: Dev only
		_, _ = w.Write(marshalJSON)
		return
	})
}

func ConfigureOpenApiFile(path string, builder *OpenApiBuilder) error {
	context, err := builder.OpenApiReflector.NewOperationContext(http.MethodGet, path)
	if err != nil {
		return err
	}
	context.AddRespStructure(new(any), func(cu *openapi.ContentUnit) {
		cu.HTTPStatus = http.StatusOK
		cu.Description = "The OpenApi file"
		cu.ContentType = "application/json"
		cu.IsDefault = true
	})
	context.AddRespStructure(new(any), func(cu *openapi.ContentUnit) {
		cu.HTTPStatus = http.StatusInternalServerError
		cu.Description = "Internal server error"
		cu.ContentType = "application/json"
	})
	err = builder.OpenApiReflector.AddOperation(context)
	if err != nil {
		return err
	}
	return nil
}
