package services

import (
	"context"
	"errors"
	"httpServer/logging"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
)

type ServiceProvider struct {
	// Logger is the logger for the service. // TODO: add scoped logger factory
	Logger logging.ILogger
	// Configuration is the Configuration for the service, some service may provide Configuration renew policy, so the value may change during the same scope.
	// Copy your own Configuration struct to use it, this is just a pointer to the Configuration in memory.
	Configuration *Configuration

	// HttpService is the ctx wrapper for the http server. It is used to run the server and handle requests.
	HttpService *HttpService

	// StoppingContext is the context which is used to stop the service. It is used to wait for the service to be stopped.
	StoppingContext context.Context
	// StoppingCancel is the cancel function for the StoppingContext. It is used to cancel the context when the service is stopped.
	StoppingCancel context.CancelFunc
}

func NewEmptyServiceProvider() *ServiceProvider {
	return &ServiceProvider{
		Logger:          nil,
		Configuration:   nil,
		HttpService:     nil,
		StoppingContext: nil,
		StoppingCancel:  nil,
	}
}

func (sp *ServiceProvider) Run(ctx context.Context) {
	sp.Logger.Information("Starting application")
	sp.StoppingContext, sp.StoppingCancel = signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	wg := sync.WaitGroup{}
	defer sp.StoppingCancel()

	go func() {
		ctx, cancel := context.WithCancel(sp.StoppingContext)
		defer cancel()
		wg.Add(1)
		defer wg.Done()
		err := sp.HttpService.Init(sp)
		if err != nil {
			sp.Logger.Warning("Error initializing http service: %v", err)
			return
		}
		err = sp.HttpService.Run(ctx)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			sp.Logger.Warning("Error running http service: %v", err)
		}
	}()

	sp.Logger.Information("Application started")
	<-sp.StoppingContext.Done()
	sp.Logger.Information("Stopping application")
	wg.Wait()
	sp.Logger.Information("Application stopped")
}

func (sp *ServiceProvider) AddHttpServiceFactory(builder func() *HttpService) {
	sp.HttpService = builder()
}
func (sp *ServiceProvider) AddHttpService(server *HttpService) {
	sp.HttpService = server
}

func (sp *ServiceProvider) AddConfigurationFactory(builder func() *Configuration) {
	sp.Configuration = builder()
}

func (sp *ServiceProvider) AddConfiguration(config *Configuration) {
	sp.Configuration = config
}

func (sp *ServiceProvider) AddLoggerFactory(builder func() logging.ILogger) {
	sp.Logger = builder()
}

func (sp *ServiceProvider) AddLogger(logger logging.ILogger) {
	sp.Logger = logger
}
