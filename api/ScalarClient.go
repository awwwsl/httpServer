package api

import (
	"net/http"
	"os"
)

func RouteScalarClient(path string, builder *RouteBuilder) {
	page, readScalarErr := os.ReadFile("assets/ScalarApiClient.html") // TODO: This is an html to cdn, use server only static files
	builder.Mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		if readScalarErr != nil {
			http.NotFound(w, r)
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/html")
		//w.Header().Set("Access-Control-Allow-Origin", "*")                   // TODO: Dev only
		//w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS") // TODO: Dev only
		w.Write(page)
	})
}
