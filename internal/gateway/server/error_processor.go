package server

import (
	"errors"
	"github.com/AskaryanKarine/bmstu-ds-4/pkg/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
)

func processError(c echo.Context, err error) error {
	var respErr models.ErrorResponse
	if errors.As(err, &respErr) {
		return c.JSON(respErr.StatusCode, respErr)
	}
	var valErr models.ValidationErrorResponse
	if errors.As(err, &valErr) {
		return c.JSON(http.StatusBadRequest, valErr)
	}
	log.Error("err = ", err)
	return c.JSON(http.StatusInternalServerError, err)
}
