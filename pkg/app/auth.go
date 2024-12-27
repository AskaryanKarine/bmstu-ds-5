package app

import (
	"fmt"
	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
)

func parseToken(token, jwksURL string) (string, error) {
	jwks, err := keyfunc.Get(jwksURL, keyfunc.Options{})
	if err != nil {
		return "", fmt.Errorf("get keyfunc: %w", err)
	}

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return jwks.Keyfunc(token)
	})

	if err != nil {
		return "", fmt.Errorf("parse jwt: %w", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid token claims type")
	}

	username, ok := claims["preferred_username"].(string)
	if !ok {
		return "", fmt.Errorf("missing username in claims")
	}

	return username, nil
}
