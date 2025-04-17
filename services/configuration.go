package services

import (
	"encoding/json"
	"httpServer/logging"
)

type EnvironmentType string

const (
	DevelopmentEnvironment EnvironmentType = "development"
	ProductionEnvironment  EnvironmentType = "production"
)

//goland:noinspection GoMixedReceiverTypes
func (env EnvironmentType) String() string {
	switch env {
	case DevelopmentEnvironment:
		return "development"
	case ProductionEnvironment:
		return "production"
	default:
		return "development" // Fallback
	}
}

//goland:noinspection GoMixedReceiverTypes
func (env EnvironmentType) MarshalJSON() ([]byte, error) {
	return []byte(`"` + env.String() + `"`), nil
}

//goland:noinspection GoMixedReceiverTypes
func (env *EnvironmentType) UnmarshalJSON(bytes []byte) error {
	var s string
	if err := json.Unmarshal(bytes, &s); err != nil {
		return err
	}
	switch s {
	case "development":
		*env = DevelopmentEnvironment
	case "production":
		*env = ProductionEnvironment
	default:
		*env = DevelopmentEnvironment // Fallback
	}
	return nil
}

type Configuration struct {
	// LogLevel is the minimum log level to be logged
	LogLevel logging.LogLevel `json:"logLevel" env:"CONFIG_LOG_LEVEL" default:"Information"`
	// Environment determines the deploy environment
	Environment EnvironmentType `json:"environment" env:"CONFIG_ENVIRONMENT" default:"production"`
	// Port is the port the server will listen on
	Port int `json:"port" env:"CONFIG_PORT" default:"8080"`
	// Host is the host the server will listen on
	Host string `json:"host" env:"CONFIG_HOST" default:"localhost"`
	// JwtSecret is the secret used to sign the JWT
	JwtSecret string `json:"jwtSecret" env:"CONFIG_JWT_SECRET" default:"YOUR_JWT_SECRET_WHICH_SHOULD_BE_LONGER_THAN_32_CHARACTERS_AND_STRONG_ENOUGH"`
	// JwtIssuer is the issuer of the JWT
	JwtIssuer string `json:"jwtIssuer" env:"CONFIG_JWT_ISSUER" default:"YOUR_JWT_ISSUER"`
}

func NewDefaultConfig() *Configuration {
	return &Configuration{
		LogLevel:    logging.Information,
		Port:        8080,
		Host:        "localhost",
		Environment: DevelopmentEnvironment,
		JwtSecret:   "YOUR_JWT_SECRET_WHICH_SHOULD_BE_LONGER_THAN_32_CHARACTERS_AND_STRONG_ENOUGH",
		JwtIssuer:   "YOUR_JWT_ISSUER",
	}
}
