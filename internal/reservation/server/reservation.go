package server

import (
	"errors"
	"github.com/AskaryanKarine/bmstu-ds-4/pkg/models"
	"github.com/AskaryanKarine/bmstu-ds-4/pkg/validation"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
	"net/http"
)

func (s *Server) createReservation(c echo.Context) error {
	username, ok := c.Get("username").(string)
	if !ok {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "failed to get username"})
	}

	var body models.ExtendedCreateReservationResponse

	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
	}

	if err := c.Validate(body); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: err.Error()})
	}

	reservationUid, err := s.rs.Create(c.Request().Context(), body, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, models.ErrorResponse{Message: "hotel uid not found"})
		}
		log.Error("err = ", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, models.ReservationResponse{ReservationUid: reservationUid})
}

func (s *Server) getReservationByUid(c echo.Context) error {
	username, ok := c.Get("username").(string)
	if !ok {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "failed to get username"})
	}

	uid := c.Param("uid")
	err := c.Validate(struct {
		Uid string `json:"uid" validate:"uuid"`
	}{uid})
	if err != nil {
		return c.JSON(http.StatusBadRequest, validation.ConvertToError(err, "failed to validate uid in path"))
	}

	res, err := s.rs.GetReservationByUUID(c.Request().Context(), uid, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, models.ErrorResponse{Message: err.Error()})
		}
		if errors.Is(err, models.WrongUsernameError) {
			return c.JSON(http.StatusForbidden, models.ErrorResponse{Message: err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, res)
}

func (s *Server) getAllReservationsByUser(c echo.Context) error {
	username, ok := c.Get("username").(string)
	if !ok {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "failed to get username"})
	}

	res, err := s.rs.GetAllReservationByUsername(c.Request().Context(), username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, models.ErrorResponse{Message: err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

func (s *Server) canceledReservations(c echo.Context) error {
	username, ok := c.Get("username").(string)
	if !ok {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "failed to get username"})
	}

	uid := c.Param("uid")
	err := c.Validate(struct {
		Uid string `json:"uid" validate:"uuid"`
	}{uid})
	if err != nil {
		return c.JSON(http.StatusBadRequest, validation.ConvertToError(err, "failed to validate uid in path"))
	}

	_, err = s.rs.GetReservationByUUID(c.Request().Context(), uid, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, models.ErrorResponse{Message: err.Error()})
		}
		if errors.Is(err, models.WrongUsernameError) {
			return c.JSON(http.StatusForbidden, models.ErrorResponse{Message: err.Error()})
		}
		log.Error("err = ", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: err.Error()})
	}

	err = s.rs.Delete(c.Request().Context(), uid)
	if err != nil {
		log.Error("err = ", err)
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusNoContent, echo.Map{})
}
