package server

import (
	"github.com/AskaryanKarine/bmstu-ds-4/pkg/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (s *Server) getLoyalty(c echo.Context) error {
	username, ok := c.Get("username").(string)
	if !ok {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{Message: "failed to get username"})
	}
	loyaltyResp, err := s.loyalty.GetLoyaltyByUser(username)
	if err != nil {
		return processError(c, err)
	}
	return c.JSON(http.StatusOK, loyaltyResp)
}
