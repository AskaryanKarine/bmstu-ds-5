package app

import (
	"github.com/AskaryanKarine/bmstu-ds-4/pkg/models"
	"github.com/AskaryanKarine/bmstu-ds-4/pkg/validation"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func SetStandardSetting(e *echo.Echo) {
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowCredentials: true,
	}))

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{}))

	e.Validator = validation.MustRegisterCustomValidator(validator.New())
}

func GetUsernameMW() func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			username := c.Request().Header.Get("X-User-Name")
			if username == "" {
				return c.JSON(http.StatusUnauthorized, models.ErrorResponse{
					Message: "failed to get X-User-Name header",
				})

			}
			c.Set("username", username)
			return next(c)
		}
	}
}
