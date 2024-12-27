package server

import (
	"github.com/AskaryanKarine/bmstu-ds-4/pkg/models"
)

type loyaltyClient interface {
	GetLoyaltyByUser(username string) (models.LoyaltyInfoResponse, error)
	DecreaseLoyalty(username string) error
	IncreaseLoyalty(username string) error
}

type reservationClient interface {
	GetHotels(page, size int, token string) (models.PaginationResponse, error)
	GetHotelByUUID(uuid, token string) (models.HotelResponse, error)
	GetReservationByUUID(username, uuid string) (models.ExtendedReservationResponse, error)
	GetReservationsByUser(username string) ([]models.ExtendedReservationResponse, error)
	CancelReservation(username, uuid string) error
	CreateReservation(model models.ExtendedCreateReservationResponse, username string) (string, error)
}

type paymentClient interface {
	GetByUUID(uuid, token string) (models.PaymentInfo, error)
	Cancel(uuid, token string) error
	CreatePayment(payment models.PaymentCreateRequest, token string) (models.ExtendedPaymentInfo, error)
}
