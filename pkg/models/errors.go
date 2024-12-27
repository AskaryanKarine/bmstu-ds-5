package models

import (
	"errors"
	"fmt"
)

type ErrorDescription struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

type ErrorResponse struct {
	// Message - информация об ошибке
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
}

func (e ErrorResponse) Error() string {
	return e.Message
}

type ValidationErrorResponse struct {
	// Message - информация об ошибке
	Message string `json:"message"`
	// Errors - Массив полей с описанием ошибки
	Errors []ErrorDescription `json:"errors"`
}

func (e ValidationErrorResponse) Error() string {
	return fmt.Sprintf("%s: %+v", e.Message, e.Errors)
}

var (
	WrongUsernameError         = errors.New("wrong username provided")
	EmptyResponseError         = errors.New("response is empty")
	UndefinedResponseCodeError = errors.New("response code undefined")
)
