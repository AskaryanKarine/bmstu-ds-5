package server

import (
	"github.com/AskaryanKarine/bmstu-ds-4/internal/gateway/clients"
	"github.com/AskaryanKarine/bmstu-ds-4/internal/gateway/config"
	"github.com/AskaryanKarine/bmstu-ds-4/pkg/app"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

const (
	defaultTimeout         = 5 * time.Second
	defaultMaxConnsPerHost = 100
)

type Server struct {
	echo        *echo.Echo
	cfg         config.Config
	loyalty     loyaltyClient
	payment     paymentClient
	reservation reservationClient
}

func NewServer(cfg config.Config) *Server {
	e := echo.New()
	client := &http.Client{
		Transport: &http.Transport{MaxConnsPerHost: defaultMaxConnsPerHost},
		Timeout:   defaultTimeout,
	}
	payment := clients.NewPaymentClient(client, cfg.PaymentService)
	reservation := clients.NewReservationClient(client, cfg.ReservationService)
	loyalty := clients.NewLoyaltyClient(client, cfg.LoyaltyService)
	s := &Server{
		echo:        e,
		cfg:         cfg,
		loyalty:     loyalty,
		payment:     payment,
		reservation: reservation,
	}

	app.SetStandardSetting(e)
	app.AddHealthCheck(e)

	api := s.echo.Group("/api/v1")

	api.GET("/hotels", s.getHotels)

	api.GET("/me", s.getUserInfo, app.GetUsernameMW())

	reservations := api.Group("/reservations", app.GetUsernameMW())
	reservations.GET("", s.getReservations, app.GetUsernameMW())
	reservations.POST("", s.createReservation, app.GetUsernameMW())
	reservations.GET("/:uid", s.getReservationsByUID, app.GetUsernameMW())
	reservations.DELETE("/:uid", s.canceledReservation, app.GetUsernameMW())

	api.GET("/loyalty", s.getLoyalty, app.GetUsernameMW())

	return s
}

func (s *Server) Run(port int) {
	app.Run(s.echo, port)
}
