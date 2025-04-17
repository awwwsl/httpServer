package main

import (
	"encoding/json"
	"errors"
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
	config := configureConfiguration(log)
	if config == nil {
		return
	}
	log = configureLogger(config.LogLevel)
	log.Information("Logging on level %s", config.LogLevel.String())
	sp := services.ServiceProvider{}
	sp.Init(config)
	httpServer := configureHttpServer(&sp)
	log.Information("Starting server on port %d", config.Port)
	err := httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Warning("Server stopped with error: %v", err)
	} else {
		log.Information("Server stopped")
	}
}

func configureLogger(logLevel logging.LogLevel) logging.ILogger {
	log := logging.NewLogger(logLevel)
	return log
}

func configureConfiguration(logger logging.ILogger) *services.Configuration {
	// Read from json
	var config services.Configuration
	configFile, err := os.ReadFile("config.json")
	if err != nil && os.IsNotExist(err) {
		// If the file does not exist, create it with default values
		config = *services.NewDefaultConfig()
		configFile, err = json.MarshalIndent(config, "", "  ")
		if err != nil {
			logger.Warning("Error creating config file: %v, exiting", err)
			return nil
		}
		err = os.WriteFile("config.json", configFile, 0644)
		if err != nil {
			logger.Warning("Error creating config file: %v, exiting", err)
			return nil
		}
		logger.Information("Config file created with default values, exiting")
		return nil
	} else {
		err = json.Unmarshal(configFile, &config)
		if err != nil {
			logger.Warning("Error reading config file: %v, creating new one", err)
			err = os.Rename("config.json", "config.json.bak")
			if err != nil {
				logger.Warning("Error backing up config file: %v, exiting", err)
				return nil
			}
			config = *services.NewDefaultConfig()
			configFile, _ = json.MarshalIndent(config, "", "  ")
			err = os.WriteFile("config.json", configFile, 0644)
			if err != nil {
				logger.Warning("Error creating config file: %v, exiting", err)
				return nil
			}
			logger.Information("Config file created with default values, exiting")
			return nil
		}
		return &config
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
