package server

import (
	"errors"
	"fmt"
	"github.com/AskaryanKarine/bmstu-ds-4/pkg/models"
	"github.com/AskaryanKarine/bmstu-ds-4/pkg/validation"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"time"
)

func (s *Server) getPaymentInfo(c echo.Context) error {
	uid := c.Param("uid")
	err := c.Validate(struct {
		Uid string `json:"uid" validate:"uuid"`
	}{uid})
	if err != nil {
		return c.JSON(http.StatusBadRequest, validation.ConvertToError(err, "failed to validate uid in path"))
	}

	res, err := s.ps.GetPaymentInfoByUUID(c.Request().Context(), uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, models.ErrorResponse{Message: fmt.Sprintf("payment info with %s uid not found", uid)})
		}
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

func (s *Server) setCanceledStatus(c echo.Context) error {
	uid := c.Param("uid")
	err := c.Validate(struct {
		Uid string `json:"uid" validate:"uuid"`
	}{uid})
	if err != nil {
		return c.JSON(http.StatusBadRequest, validation.ConvertToError(err, "failed to validate uid in path"))
	}
	err = s.ps.Delete(c.Request().Context(), uid)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusNoContent, echo.Map{})
}

func (s *Server) CreatePayment(c echo.Context) error {
	var body models.PaymentCreateRequest
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest,
			models.ErrorResponse{Message: err.Error()},
		)
	}

	if err := c.Validate(body); err != nil {
		return c.JSON(http.StatusBadRequest, validation.ConvertToError(err, "failed to validate body"))
	}

	startDate, err := time.Parse("2006-01-02", body.StartDate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
	}
	endDate, err := time.Parse("2006-01-02", body.EndDate)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
	}

	days := endDate.Sub(startDate).Hours() / 24
	cost := float64(body.Price) * days
	costWithDiscount := cost - (cost * (float64(body.Discount) * 0.01))

	payment := models.PaymentInfo{
		Status: models.PAID,
		Price:  int(costWithDiscount),
	}
	paymentUid, err := s.ps.Create(c.Request().Context(), payment)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, models.ExtendedPaymentInfo{PaymentUid: paymentUid, PaymentInfo: payment})
}
