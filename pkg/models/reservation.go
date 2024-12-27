package models

type ReservationResponse struct {
	ReservationUid string            `json:"reservationUid"`
	Hotel          HotelInfo         `json:"hotel"`
	StartDate      string            `json:"startDate"`
	EndDate        string            `json:"endDate"`
	Status         PaymentStatusType `json:"status"`
	Payment        PaymentInfo       `json:"payment"`
}

type CreateReservationRequest struct {
	HotelUid  string `json:"hotelUid" validate:"required,uuid"`
	StartDate string `json:"startDate" validate:"required,IsISO8601"`
	EndDate   string `json:"endDate" validate:"required,IsISO8601"`
}

type CreateReservationResponse struct {
	ReservationUid string            `json:"reservationUid"`
	Discount       int               `json:"discount"`
	Status         PaymentStatusType `json:"status"`
	Payment        PaymentInfo       `json:"payment"`
	CreateReservationRequest
}

type ExtendedReservationResponse struct {
	ReservationResponse
	PaymentUID string `json:"paymentUid"`
}

type ExtendedCreateReservationResponse struct {
	CreateReservationRequest
	PaymentUid string `json:"paymentUid" validate:"required,uuid"`
}
