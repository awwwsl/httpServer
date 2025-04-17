package services

import (
	"github.com/patrickmn/go-cache"
	"httpServer/logging"
)

type ServiceProvider struct {
	Logger             logging.ILogger
	AllowedCredentials *cache.Cache
	// Configuration is the Configuration for the service, some service may provide Configuration renew policy, so the value may change during the same scope.
	// Copy your own Configuration struct to use it, this is just a pointer to the Configuration in memory.
	Configuration    *Configuration
	AuthorizeService IAuthorizeService
}

func (sp *ServiceProvider) Init(configuration *Configuration) {
	sp.Configuration = configuration
	sp.UseDefaultLogger(configuration)
	sp.UseDefaultAuthorizeService()
}

func (sp *ServiceProvider) UseDefaultLogger(configuration *Configuration) {
	sp.Logger = logging.NewLogger(configuration.LogLevel)
}

func (sp *ServiceProvider) UseDefaultAuthorizeService() {
	sp.AuthorizeService = &authorizeService{
		serviceProvider: sp,
	}
}
