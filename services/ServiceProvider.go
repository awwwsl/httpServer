package services

import (
	"context"
	"github.com/patrickmn/go-cache"
	"httpServer/logging"
	"os/signal"
	"syscall"
)

type ServiceProvider struct {
	Logger             logging.ILogger
	AllowedCredentials *cache.Cache
	// Configuration is the Configuration for the service, some service may provide Configuration renew policy, so the value may change during the same scope.
	// Copy your own Configuration struct to use it, this is just a pointer to the Configuration in memory.
	Configuration    *Configuration
	AuthorizeService IAuthorizeService
	StoppingContext  context.Context
	StoppingCancel   context.CancelFunc
}

func (sp *ServiceProvider) Init(configuration *Configuration) {
	sp.Configuration = configuration
	sp.UseDefaultLogger(configuration)
	sp.UseDefaultAuthorizeService()
	sp.StoppingContext, sp.StoppingCancel = signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sp.StoppingContext.Done()
		defer sp.StoppingCancel() // someone says without this will cause leak, idk if this is true but added anyway
		sp.Logger.Information("Shutting down application")
	}()
}

func (sp *ServiceProvider) UseDefaultLogger(configuration *Configuration) {
	sp.Logger = logging.NewLogger(configuration.LogLevel)
}

func (sp *ServiceProvider) UseDefaultAuthorizeService() {
	sp.AuthorizeService = &authorizeService{
		serviceProvider: sp,
	}
}
