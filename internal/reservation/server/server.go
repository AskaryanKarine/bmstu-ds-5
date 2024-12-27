package server

import (
	"context"
	innermodels "github.com/AskaryanKarine/bmstu-ds-4/internal/reservation/models"
	"github.com/AskaryanKarine/bmstu-ds-4/pkg/app"
	"github.com/AskaryanKarine/bmstu-ds-4/pkg/models"
	"github.com/labstack/echo/v4"
)

type Server struct {
	echo *echo.Echo
	hs   hotelStorage
	rs   reservationStorage
}

//go:generate minimock -o mocks_storage.go -i hotelStorage -g
type hotelStorage interface {
	GetAllHotels(ctx context.Context, pagination innermodels.PaginationParams) ([]models.HotelResponse, int, error)
	GetHotelInfoByUUID(ctx context.Context, uuid string) (models.HotelResponse, error)
}

type reservationStorage interface {
	GetReservationByUUID(ctx context.Context, uuid, username string) (models.ExtendedReservationResponse, error)
	GetAllReservationByUsername(ctx context.Context, username string) ([]models.ExtendedReservationResponse, error)
	Delete(ctx context.Context, uuid string) error
	Create(ctx context.Context, reservation models.ExtendedCreateReservationResponse, username string) (string, error)
}

func New(hs hotelStorage, rs reservationStorage) *Server {
	e := echo.New()
	s := &Server{
		echo: e,
		hs:   hs,
		rs:   rs,
	}

	app.SetStandardSetting(e)
	app.AddHealthCheck(e)

	api := s.echo.Group("/api/v1")

	api.GET("/hotels", s.getAllHotels)
	api.GET("/hotels/:uid", s.getHotelByUID)

	reservations := api.Group("/reservations")
	reservations.GET("", s.getAllReservationsByUser, app.GetUsernameMW())
	reservations.POST("", s.createReservation, app.GetUsernameMW())
	reservations.GET("/:uid", s.getReservationByUid, app.GetUsernameMW())
	reservations.DELETE("/:uid", s.canceledReservations, app.GetUsernameMW())

	return s
}

func (s *Server) Run(port int) {
	app.Run(s.echo, port)
}
