package server

import (
	"github.com/AskaryanKarine/bmstu-ds-4/pkg/models"
	"github.com/AskaryanKarine/bmstu-ds-4/pkg/validation"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

func (s *Server) getHotels(c echo.Context) error {
	token, ok := c.Get("token").(string)
	if !ok {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "failed to get token"})
	}
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "page must be integer"})
	}
	size, err := strconv.Atoi(c.QueryParam("size"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "size must be integer"})
	}
	resp, err := s.reservation.GetHotels(page, size, token)
	if err != nil {
		return processError(c, err)
	}
	return c.JSON(http.StatusOK, resp)
}

func (s *Server) getReservations(c echo.Context) error {
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

	return c.JSON(http.StatusOK, reservationsResp)
}

func (s *Server) getReservationsByUID(c echo.Context) error {
	uid := c.Param("uid")
	err := c.Validate(struct {
		Uid string `json:"uid" validate:"uuid"`
	}{uid})
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
	}

	token, ok := c.Get("token").(string)
	if !ok {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "failed to get token"})
	}

	reservationsByUser, err := s.reservation.GetReservationByUUID(token, uid)
	if err != nil {
		return processError(c, err)
	}

	paymentInfo, err := s.payment.GetByUUID(reservationsByUser.PaymentUID, token)
	if err != nil {
		return processError(c, err)
	}
	reservationsByUser.Payment = paymentInfo
	return c.JSON(http.StatusOK, reservationsByUser.ReservationResponse)
}

func (s *Server) canceledReservation(c echo.Context) error {
	uid := c.Param("uid")
	err := c.Validate(struct {
		Uid string `json:"uid" validate:"uuid"`
	}{uid})
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
	}

	token, ok := c.Get("token").(string)
	if !ok {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "failed to get token"})
	}

	err = s.reservation.CancelReservation(token, uid)
	if err != nil {
		return processError(c, err)
	}

	reservations, err := s.reservation.GetReservationByUUID(token, uid)
	if err != nil {
		return processError(c, err)
	}

	err = s.payment.Cancel(reservations.PaymentUID, token)
	if err != nil {
		return processError(c, err)
	}

	err = s.loyalty.DecreaseLoyalty(token)
	if err != nil {
		return processError(c, err)
	}

	return c.JSON(http.StatusNoContent, echo.Map{})
}

func (s *Server) createReservation(c echo.Context) error {
	token, ok := c.Get("token").(string)
	if !ok {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "failed to get token"})
	}

	var body models.CreateReservationRequest
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
	}

	if err := c.Validate(body); err != nil {
		return c.JSON(http.StatusBadRequest, validation.ConvertToError(err, "failed to validate body request"))
	}

	hotelInfo, err := s.reservation.GetHotelByUUID(body.HotelUid, token) // отсюда только price
	if err != nil {
		return processError(c, err)
	}

	loyalty, err := s.loyalty.GetLoyaltyByUser(token) // отсюда только discount
	if err != nil {
		return processError(c, err)
	}

	// передавать в P инфу о discount, price, start date, end date
	extendedPaymentInfo, err := s.payment.CreatePayment(models.PaymentCreateRequest{
		Price:     hotelInfo.Price,
		Discount:  loyalty.Discount,
		StartDate: body.StartDate,
		EndDate:   body.EndDate,
	}, token)

	if err != nil {
		return processError(c, err)
	}

	err = s.loyalty.IncreaseLoyalty(token)
	if err != nil {
		return processError(c, err)
	}

	reservationUid, err := s.reservation.CreateReservation(models.ExtendedCreateReservationResponse{
		CreateReservationRequest: body,
		PaymentUid:               extendedPaymentInfo.PaymentUid,
	}, token)

	if err != nil {
		return processError(c, err)
	}

	reservation, err := s.reservation.GetReservationByUUID(token, reservationUid)
	if err != nil {
		return processError(c, err)
	}
	reservation.Payment = extendedPaymentInfo.PaymentInfo
	response := models.CreateReservationResponse{
		ReservationUid:           reservationUid,
		Discount:                 loyalty.Discount,
		Status:                   reservation.Status,
		Payment:                  extendedPaymentInfo.PaymentInfo,
		CreateReservationRequest: body,
	}
	return c.JSON(http.StatusOK, response)
}
