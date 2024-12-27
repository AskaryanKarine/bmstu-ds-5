package server

import (
	"errors"
	innermodels "github.com/AskaryanKarine/bmstu-ds-4/internal/reservation/models"
	"github.com/AskaryanKarine/bmstu-ds-4/pkg/models"
	"github.com/AskaryanKarine/bmstu-ds-4/pkg/validation"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
)

func (s *Server) getAllHotels(c echo.Context) error {
	var qParams innermodels.PaginationParams

	if err := c.Bind(&qParams); err != nil {
		return c.JSON(http.StatusBadRequest,
			models.ErrorResponse{Message: err.Error()},
		)
	}

	if err := c.Validate(qParams); err != nil {
		return c.JSON(http.StatusBadRequest, validation.ConvertToError(err, "failed to validate query params"))
	}

	result, count, err := s.hs.GetAllHotels(c.Request().Context(), qParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, models.PaginationResponse{
		Page:          qParams.Page,
		PageSize:      qParams.Size,
		TotalElements: count,
		Items:         result,
	})
}

func (s *Server) getHotelByUID(c echo.Context) error {
	uid := c.Param("uid")
	err := c.Validate(struct {
		Uid string `json:"uid" validate:"uuid"`
	}{uid})
	if err != nil {
		return c.JSON(http.StatusBadRequest, validation.ConvertToError(err, "failed to validate uid in path"))
	}

	hotelInfo, err := s.hs.GetHotelInfoByUUID(c.Request().Context(), uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, models.ErrorResponse{Message: err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, hotelInfo)
}
