package server

import (
	"context"
	"github.com/AskaryanKarine/bmstu-ds-4/pkg/app"
	"github.com/AskaryanKarine/bmstu-ds-4/pkg/config"
	"github.com/AskaryanKarine/bmstu-ds-4/pkg/models"
	"github.com/labstack/echo/v4"
)

type Server struct {
	echo *echo.Echo
	lr   loyaltyRepository
}

//go:generate minimock -o mocks_storage.go -g
type loyaltyRepository interface {
	GetByUser(ctx context.Context, username string) (models.LoyaltyInfoResponse, error)
	UpdateByUser(ctx context.Context, username string, usersLoyalty models.LoyaltyInfoResponse) error
}

func NewServer(lr loyaltyRepository, cfg config.Config) *Server {
	e := echo.New()
	s := &Server{
		echo: e,
		lr:   lr,
	}

	app.SetStandardSetting(e)
	app.AddHealthCheck(e)

	api := s.echo.Group("/api/v1", app.GetUsernameMW(cfg.JWKsURl))

	api.GET("/loyalty", s.getLoyaltyByUser)

	reservations := api.Group("/reservations")
	reservations.DELETE("/decrease", s.decreaseCounter)
	reservations.POST("/increase", s.increaseCounter)

	return s
}

func (s *Server) Run(port int) {
	app.Run(s.echo, port)
}
