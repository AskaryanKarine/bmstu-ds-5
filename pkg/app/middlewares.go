package app

import (
	"context"
	"fmt"
	"github.com/AskaryanKarine/bmstu-ds-4/pkg/validation"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"strings"
)

func SetStandardSetting(e *echo.Echo) {
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowCredentials: true,
	}))

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{}))

	e.Validator = validation.MustRegisterCustomValidator(validator.New())
}

func GetUsernameMW(jwksURL string) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			header := c.Request().Header.Get("Authorization")
			if header == "" {
				return c.NoContent(http.StatusUnauthorized)
			}

			prefix := "Bearer "
			if !strings.HasPrefix(header, prefix) {
				return c.NoContent(http.StatusUnauthorized)
			}

			token := strings.TrimPrefix(header, prefix)

			username, err := parseToken(token, jwksURL)
			fmt.Println(username, err)
			if err != nil {
				return c.NoContent(http.StatusUnauthorized)
			}

			c.Set("username", username)
			c.Set("token", token)
			return next(c)
		}
	}
}

func GetToken(ctx context.Context) string {
	return ""
}
