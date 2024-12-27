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
	GetHotels(page, size int) (models.PaginationResponse, error)
	GetHotelByUUID(uuid string) (models.HotelResponse, error)
	GetReservationByUUID(username, uuid string) (models.ExtendedReservationResponse, error)
	GetReservationsByUser(username string) ([]models.ExtendedReservationResponse, error)
	CancelReservation(username, uuid string) error
	CreateReservation(model models.ExtendedCreateReservationResponse, username string) (string, error)
}

type paymentClient interface {
	GetByUUID(uuid string) (models.PaymentInfo, error)
	Cancel(uuid string) error
	CreatePayment(payment models.PaymentCreateRequest) (models.ExtendedPaymentInfo, error)
}
