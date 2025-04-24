package main

import (
	"context"
	"encoding/json"
	"github.com/swaggest/openapi-go/openapi3"
	"httpServer/api"
	"httpServer/logging"
	"httpServer/services"
	"net/http"
	"os"
	"strconv"
)

func main() {
	log := configureLogger(logging.Information)
	config, err := configureConfiguration(log)
	if err != nil {
		log.Warning("Error configuring application: %v", err)
		return
	}
	log = configureLogger(config.LogLevel)
	log.Information("Logging on level %s", config.LogLevel.String())

	sp := services.NewEmptyServiceProvider()
	sp.AddLogger(log)
	sp.AddConfiguration(config)
	sp.AddHttpService(func() *services.HttpService {
		server := configureHttpServer(sp)
		return services.NewHttpService(server, nil)
	}())
	sp.Run(context.Background())

	log.Information("Exiting")
}

func configureLogger(logLevel logging.LogLevel) logging.ILogger {
	log := logging.NewLogger(logLevel)
	return log
}

func configureConfiguration(logger logging.ILogger) (*services.Configuration, error) {
	// Read from json
	var config services.Configuration
	configFile, err := os.ReadFile("config.json")
	if err != nil && os.IsNotExist(err) {
		// If the file does not exist, create it with default values
		config = *services.NewDefaultConfig()
		configFile, err = json.MarshalIndent(config, "", "  ")
		if err != nil {
			logger.Warning("Error creating config file: %v, exiting", err)
			return nil, err
		}
		err = os.WriteFile("config.json", configFile, 0644)
		if err != nil {
			logger.Warning("Error creating config file: %v, exiting", err)
			return nil, err
		}
		logger.Information("Config file created with default values, exiting")
		return nil, err
	} else {
		err = json.Unmarshal(configFile, &config)
		if err != nil {
			logger.Warning("Error reading config file: %v, creating new one", err)
			err = os.Rename("config.json", "config.json.bak")
			if err != nil {
				logger.Warning("Error backing up config file: %v, exiting", err)
				return nil, err
			}
			config = *services.NewDefaultConfig()
			configFile, _ = json.MarshalIndent(config, "", "  ")
			err = os.WriteFile("config.json", configFile, 0644)
			if err != nil {
				logger.Warning("Error creating config file: %v, exiting", err)
				return nil, err
			}
			logger.Information("Config file created with default values, exiting")
			return nil, err
		}
		return &config, nil
	}
	// Read from env
}

func configureHttpServer(sp *services.ServiceProvider) *http.Server {
	var err error
	routeBuilder := api.NewRouteBuilder(sp)
	openApiBuilder := api.NewOpenApiBuilder()
	configureOpenApiBasics(openApiBuilder.OpenApiReflector)
	routeBuilder.Mux.Handle("/", http.RedirectHandler("/api/openapi", http.StatusFound))
	api.RouteScalarClient("/api/openapi", routeBuilder)
	api.RouteOpenApiFile("/api/openapi/openapi.json", routeBuilder, openApiBuilder)
	err = api.ConfigureOpenApiFile("/api/openapi/openapi.json", openApiBuilder)
	if err != nil {
		sp.Logger.Warning("Error configuring OpenApi file: %v", err)
		err = nil
	}
	api.RoutePerlinNoise("/api/perlin_noise", routeBuilder)
	err = api.ConfigurePerlinNoise("/api/perlin_noise", openApiBuilder)
	if err != nil {
		sp.Logger.Warning("Error configuring Perlin noise: %v", err)
		err = nil
	}
	api.RouteDrunkBishop("/api/drunk_bishop", routeBuilder)
	err = api.ConfigureDrunkBishop("/api/drunk_bishop", openApiBuilder)
	if err != nil {
		sp.Logger.Warning("Error configuring Drunk Bishop: %v", err)
		err = nil
	}
	api.RouteBrainFxxkInterpretor("/api/brain_fxxk_interpretor", routeBuilder)
	err = api.ConfigureBrainFxxkInterpretor("/api/brain_fxxk_interpretor", openApiBuilder)
	if err != nil {
		sp.Logger.Warning("Error configuring BrainFxxk Interpretor: %v", err)
		err = nil
	}

	if true {
		sp.Logger.Warning("Exposing pprof at /api/pprof, this is not recommended in production")
		api.RoutePProf("/api/pprof", routeBuilder)
		err = api.ConfigurePProf("/api/pprof", openApiBuilder)
		if err != nil {
			sp.Logger.Warning("Error configuring PProf: %v", err)
			err = nil
		}
	}
	server := http.Server{
		Addr:    ":" + strconv.Itoa(sp.Configuration.Port),
		Handler: newLoggingServeMux(sp.Logger, routeBuilder.Mux),
	}
	return &server
}

func configureOpenApiBasics(reflector *openapi3.Reflector) {
	reflector.Spec = &openapi3.Spec{
		Openapi: "3.0.4",
	}
	reflector.Spec.Info.
		WithTitle("awwwsl Go HttpServer").
		WithVersion("0.1.0").
		WithDescription("foobarbaz")
}
