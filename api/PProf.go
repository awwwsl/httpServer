package api

import (
	"github.com/swaggest/openapi-go"
	"net/http"
	httpPProf "net/http/pprof"
	runtimePProf "runtime/pprof"
	"strconv"
	"strings"
)

func RoutePProf(path string, builder *RouteBuilder) {
	//builder.Mux.HandleFunc(path, pprof.Index) // this index uses hard-coded route path so it will not work
	builder.Mux.HandleFunc(path+"/cmdline", httpPProf.Cmdline)
	builder.Mux.HandleFunc(path+"/profile", httpPProf.Profile)
	builder.Mux.HandleFunc(path+"/symbol", httpPProf.Symbol)
	builder.Mux.HandleFunc(path+"/trace", httpPProf.Trace)
	// Replaced with dynamic resolved runtime/pprof.Profile
	//builder.Mux.HandleFunc(path+"/allocs", pprof.Handler("allocs").ServeHTTP)
	//builder.Mux.HandleFunc(path+"/block", pprof.Handler("block").ServeHTTP)
	//builder.Mux.HandleFunc(path+"/goroutine", pprof.Handler("goroutine").ServeHTTP)
	//builder.Mux.HandleFunc(path+"/heap", pprof.Handler("heap").ServeHTTP)
	//builder.Mux.HandleFunc(path+"/mutex", pprof.Handler("mutex").ServeHTTP)
	//builder.Mux.HandleFunc(path+"/threadcreate", pprof.Handler("threadcreate").ServeHTTP)

	profiles := runtimePProf.Profiles()
	for _, profile := range profiles {
		builder.Mux.HandleFunc(path+"/"+profile.Name(), func(w http.ResponseWriter, r *http.Request) {
			debugStr := r.URL.Query().Get("debug")
			debug := 0
			if debugStr != "" {
				var err error
				debug, err = strconv.Atoi(debugStr)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					_, _ = w.Write([]byte("Invalid debug parameter"))
					return
				}
			}
			_ = profile.WriteTo(w, debug)
		})
	}
}

func ConfigurePProf(path string, builder *OpenApiBuilder) error {
	var err error
	var context openapi.OperationContext
	// cmdline
	context, err = builder.OpenApiReflector.NewOperationContext(http.MethodGet, path+"/cmdline")
	if err != nil {
		return err
	}
	context.SetTags("debug")
	context.SetSummary("cmdline")
	context.SetDescription(`Cmdline responds with the running program's command line, with arguments separated by NUL bytes. The package initialization registers it as /debug/pprof/cmdline.`)
	context.AddRespStructure(new(string), func(cu *openapi.ContentUnit) {
		cu.HTTPStatus = http.StatusOK
		cu.ContentType = `text/plain`
		cu.Description = `Cmdline responds with the running program's command line, with arguments separated by NUL bytes. The package initialization registers it as /debug/pprof/cmdline.`
	})
	err = builder.OpenApiReflector.AddOperation(context)
	if err != nil {
		return err
	}

	context, err = builder.OpenApiReflector.NewOperationContext(http.MethodGet, path+"/profile")
	if err != nil {
		return err
	}
	context.SetTags("debug")
	context.SetSummary("profile")
	context.SetDescription(`Profile responds with the pprof-formatted cpu profile. Profiling lasts for duration specified in seconds GET parameter, or for 30 seconds if not specified. The package initialization registers it as /debug/pprof/profile.`)
	context.AddRespStructure(new(string), func(cu *openapi.ContentUnit) {
		cu.HTTPStatus = http.StatusOK
		cu.ContentType = `application/octet-stream`
		cu.Description = `Profile responds with the pprof-formatted cpu profile. Profiling lasts for duration specified in seconds GET parameter, or for 30 seconds if not specified. The package initialization registers it as /debug/pprof/profile.`
	})
	err = builder.OpenApiReflector.AddOperation(context)
	if err != nil {
		return err
	}

	context, err = builder.OpenApiReflector.NewOperationContext(http.MethodGet, path+"/symbol")
	if err != nil {
		return err
	}
	context.SetTags("debug")
	context.SetSummary("symbol")
	context.SetDescription(`Symbol looks up the program counters listed in the request, responding with a table mapping program counters to function names. The package initialization registers it as /debug/pprof/symbol.`)
	context.AddRespStructure(new(string), func(cu *openapi.ContentUnit) {
		cu.HTTPStatus = http.StatusOK
		cu.ContentType = `text/plain`
		cu.Description = `Symbol looks up the program counters listed in the request, responding with a table mapping program counters to function names. The package initialization registers it as /debug/pprof/symbol.`
	})
	err = builder.OpenApiReflector.AddOperation(context)
	if err != nil {
		return err
	}

	context, err = builder.OpenApiReflector.NewOperationContext(http.MethodGet, path+"/trace")
	if err != nil {
		return err
	}
	context.SetTags("debug")
	context.SetSummary("trace")
	context.SetDescription(`Trace responds with the execution trace in binary form. Tracing lasts for duration specified in seconds GET parameter, or for 1 second if not specified. The package initialization registers it as /debug/pprof/trace.`)
	context.AddRespStructure(new(string), func(cu *openapi.ContentUnit) {
		cu.HTTPStatus = http.StatusOK
		cu.ContentType = `application/octet-stream`
		cu.Description = `Trace responds with the execution trace in binary form. Tracing lasts for duration specified in seconds GET parameter, or for 1 second if not specified. The package initialization registers it as /debug/pprof/trace.`
	})
	err = builder.OpenApiReflector.AddOperation(context)
	if err != nil {
		return err
	}

	profiles := runtimePProf.Profiles()
	for _, profile := range profiles {
		context, err = builder.OpenApiReflector.NewOperationContext(http.MethodGet, path+"/"+profile.Name())
		if err != nil {
			return err
		}
		context.SetTags("debug")
		context.SetSummary(profile.Name())
		context.SetDescription(`Count ` + strconv.Itoa(profile.Count()) + "\n" + getDescriptionOrDefault(profile.Name()))
		context.AddRespStructure(new(string), func(cu *openapi.ContentUnit) {
			cu.HTTPStatus = http.StatusOK
			cu.ContentType = `application/octet-stream`
			cu.Description = getDescriptionOrDefault(profile.Name())
		})
		err = builder.OpenApiReflector.AddOperation(context)
		if err != nil {
			return err
		}
	}
	return nil
}

func getDescriptionOrDefault(name string) string {
	//goroutine    - stack traces of all current goroutines
	//heap         - a sampling of memory allocations of live objects
	//allocs       - a sampling of all past memory allocations
	//threadcreate - stack traces that led to the creation of new OS threads
	//block        - stack traces that led to blocking on synchronization primitives
	//mutex        - stack traces of holders of contended mutexes
	name = strings.ToLower(name)
	switch name {
	case "goroutine":
		return "stack traces of all current goroutines"
	case "heap":
		return "a sampling of memory allocations of live objects"
	case "allocs":
		return "a sampling of all past memory allocations"
	case "threadcreate":
		return "stack traces that led to the creation of new OS threads"
	case "block":
		return "stack traces that led to blocking on synchronization primitives"
	case "mutex":
		return "stack traces of holders of contended mutexes"
	default:
		return name
	}
}
