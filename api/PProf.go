package api

import (
	"github.com/swaggest/openapi-go"
	"net/http"
	"net/http/pprof"
)

func RoutePProf(path string, builder *RouteBuilder) {
	//builder.Mux.HandleFunc(path, pprof.Index) // this index uses hard-coded route path so it will not work
	builder.Mux.HandleFunc(path+"/cmdline", pprof.Cmdline)
	builder.Mux.HandleFunc(path+"/profile", pprof.Profile)
	builder.Mux.HandleFunc(path+"/symbol", pprof.Symbol)
	builder.Mux.HandleFunc(path+"/trace", pprof.Trace)
	builder.Mux.HandleFunc(path+"/allocs", pprof.Handler("allocs").ServeHTTP)
	builder.Mux.HandleFunc(path+"/block", pprof.Handler("block").ServeHTTP)
	builder.Mux.HandleFunc(path+"/goroutine", pprof.Handler("goroutine").ServeHTTP)
	builder.Mux.HandleFunc(path+"/heap", pprof.Handler("heap").ServeHTTP)
	builder.Mux.HandleFunc(path+"/mutex", pprof.Handler("mutex").ServeHTTP)
	builder.Mux.HandleFunc(path+"/threadcreate", pprof.Handler("threadcreate").ServeHTTP)
}

func ConfigurePProf(path string, builder *OpenApiBuilder) error {
	var err error
	// this index uses hard-coded route path so it will not work
	// err := configureInnerPProf(path, "text/html", "", "pprof index page", builder)
	// if err != nil {
	// 	return err
	// }
	err = configureInnerPProf(path, "text/plain", "/cmdline", "pprof command line", builder)
	if err != nil {
		return err
	}
	err = configureInnerPProf(path, "application/octet-stream", "/profile", "pprof profile", builder)
	if err != nil {
		return err
	}
	err = configureInnerPProf(path, "text/plain", "/symbol", "pprof symbol", builder)
	if err != nil {
		return err
	}
	err = configureInnerPProf(path, "text/plain", "/trace", "pprof trace", builder)
	if err != nil {
		return err
	}
	err = configureInnerPProf(path, "application/octet-stream", "/allocs", "pprof allocs", builder)
	if err != nil {
		return err
	}
	err = configureInnerPProf(path, "application/octet-stream", "/block", "pprof block", builder)
	if err != nil {
		return err
	}
	err = configureInnerPProf(path, "application/octet-stream", "/goroutine", "pprof goroutine", builder)
	if err != nil {
		return err
	}
	err = configureInnerPProf(path, "application/octet-stream", "/heap", "pprof heap", builder)
	if err != nil {
		return err
	}
	err = configureInnerPProf(path, "application/octet-stream", "/mutex", "pprof mutex", builder)
	if err != nil {
		return err
	}
	err = configureInnerPProf(path, "application/octet-stream", "/threadcreate", "pprof threadcreate", builder)
	if err != nil {
		return err
	}
	return nil
}

func configureInnerPProf(path, contentType, suffix, description string, builder *OpenApiBuilder) error {
	context, err := builder.OpenApiReflector.NewOperationContext(http.MethodGet, path+suffix)
	if err != nil {
		return err
	}
	context.SetSummary(description)
	context.SetDescription(description)
	context.SetTags("debug")
	context.AddRespStructure(new(string), func(cu *openapi.ContentUnit) {
		cu.HTTPStatus = http.StatusOK
		cu.Description = description
		cu.ContentType = contentType
		cu.IsDefault = true
	})
	err = builder.OpenApiReflector.AddOperation(context)
	if err != nil {
		return err
	}
	return nil
}
