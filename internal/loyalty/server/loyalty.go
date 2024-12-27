package server

import (
	"errors"
	"fmt"
	"github.com/AskaryanKarine/bmstu-ds-4/pkg/models"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
)

var statusMap = map[models.LoyaltyStatusType]struct {
	discount, count int
	prev, next      models.LoyaltyStatusType
}{
	models.BRONZE: {5, 0, models.BRONZE, models.SILVER},
	models.SILVER: {7, 10, models.BRONZE, models.GOLD},
	models.GOLD:   {10, 20, models.SILVER, models.GOLD},
}

func (s *Server) getLoyaltyByUser(c echo.Context) error {
	username, ok := c.Get("username").(string)
	if !ok {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "failed to get username"})
	}

	res, err := s.lr.GetByUser(c.Request().Context(), username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, models.ErrorResponse{Message: fmt.Sprintf("user %s not found", username)})
		}
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}

func (s *Server) decreaseCounter(c echo.Context) error {
	username, ok := c.Get("username").(string)
	if !ok {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "failed to get username"})
	}

	res, err := s.lr.GetByUser(c.Request().Context(), username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, models.ErrorResponse{Message: fmt.Sprintf("user %s not found", username)})
		}
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: err.Error()})
	}
	*res.ReservationCount--
	if *res.ReservationCount < statusMap[res.Status].count {
		res.Status = statusMap[res.Status].prev
		res.Discount = statusMap[res.Status].discount
	}
	err = s.lr.UpdateByUser(c.Request().Context(), username, res)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusNoContent, echo.Map{})
}

func (s *Server) increaseCounter(c echo.Context) error {
	username, ok := c.Get("username").(string)
	if !ok {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "failed to get username"})
	}
	res, err := s.lr.GetByUser(c.Request().Context(), username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, models.ErrorResponse{Message: fmt.Sprintf("user %s not found", username)})
		}
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: err.Error()})
	}
	*res.ReservationCount++
	if *res.ReservationCount > statusMap[res.Status].count {
		res.Status = statusMap[res.Status].next
		res.Discount = statusMap[res.Status].discount
	}
	err = s.lr.UpdateByUser(c.Request().Context(), username, res)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusNoContent, echo.Map{})
}
