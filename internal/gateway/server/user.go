package server

import (
	"github.com/AskaryanKarine/bmstu-ds-4/pkg/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (s *Server) getUserInfo(c echo.Context) error {
	token, ok := c.Get("token").(string)
	if !ok {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "failed to get token"})
	}

	reservationsByUser, err := s.reservation.GetReservationsByUser(token)
	if err != nil {
		return processError(c, err)
	}

	reservationsResp := make([]models.ReservationResponse, 0, len(reservationsByUser))
	for i := range reservationsByUser {
		paymentInfo, err := s.payment.GetByUUID(reservationsByUser[i].PaymentUID, token)
		if err != nil {
			return processError(c, err)
		}
		reservationsByUser[i].Payment = paymentInfo
		reservationsResp = append(reservationsResp, reservationsByUser[i].ReservationResponse)
	}
	loyaltyResp, err := s.loyalty.GetLoyaltyByUser(token)
	if err != nil {
		return processError(c, err)
	}
	loyaltyResp.ReservationCount = nil
	return c.JSON(http.StatusOK, echo.Map{
		"reservations": reservationsResp,
		"loyalty":      loyaltyResp,
	})
}
