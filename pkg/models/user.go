package models

type UserInfoResponse struct {
	Reservations []ReservationResponse `json:"reservations"`
	Loyalty      LoyaltyInfoResponse   `json:"loyalty"`
}
