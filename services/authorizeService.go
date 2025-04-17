package services

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type authorizeService struct {
	serviceProvider *ServiceProvider
}

func (j authorizeService) addCredentials(credential string) error {
	//TODO implement me
	panic("implement me")
}

func (j authorizeService) verifyCredentials(credential string) (string, error) {
	//TODO implement me
	panic("implement me")
}

type IAuthorizeService interface {
	addCredentials(credential string) error
	verifyCredentials(credential string) (string, error)
}

func (j authorizeService) generateJwt(claimsBuilder func(claim jwt.Claims)) (string, error) {
	j.serviceProvider.Logger.Verbose("Generating new JWT")
	secret := j.serviceProvider.Configuration.JwtSecret
	config := j.serviceProvider.Configuration
	claims := jwt.MapClaims{
		"iss": config.JwtIssuer,
		"exp": jwt.NewNumericDate(time.Now().Add(time.Minute * 10)),
		"iat": jwt.NewNumericDate(time.Now()),
	}
	claimsBuilder(claims)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	j.serviceProvider.Logger.Trace("Generated new JWT \"%s\" with secret \"%s\"", tokenString, secret)
	return tokenString, nil
}

func (j authorizeService) validateJwt(token string) (jwt.Claims, error) {
	logger := j.serviceProvider.Logger
	config := j.serviceProvider.Configuration
	logger.Verbose("Validating JWT")
	jwt, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logger.Warning("Unexpected signing method: %v", token.Header["alg"])
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(config.JwtSecret), nil
	},
		jwt.WithIssuer(config.JwtIssuer),
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
	)
	if err != nil {
		logger.Verbose("Error parsing jwt \"%s\" with secret \"%s\"\nparse options: %s", err, j.serviceProvider.Configuration.JwtSecret, token)
		return nil, err
	}
	return jwt.Claims, nil
}
