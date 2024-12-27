package server

import (
	"github.com/AskaryanKarine/bmstu-ds-4/pkg/validation"
	"github.com/go-playground/validator/v10"
	"github.com/gojuno/minimock/v3"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_setCanceledStatus(t *testing.T) {
	mc := minimock.NewController(t)
	e := echo.New()
	e.Validator = validation.MustRegisterCustomValidator(validator.New())
	type fields struct {
		echo *echo.Echo
		ps   paymentStorage
	}
	tests := []struct {
		name               string
		fields             fields
		expectedHTTPStatus int
		pathParams         string
	}{
		{
			name: "http-204: ok",
			fields: fields{
				echo: e,
				ps:   NewPaymentStorageMock(mc).DeleteMock.Return(nil),
			},
			expectedHTTPStatus: http.StatusNoContent,
			pathParams:         "f440db59-70df-4d8f-aa36-99d85f3ca589",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				echo: tt.fields.echo,
				ps:   tt.fields.ps,
			}
			r := httptest.NewRequest(http.MethodGet, "/test", nil)
			w := httptest.NewRecorder()
			c := s.echo.NewContext(r, w)
			c.SetPath("/:uid")
			c.SetParamNames("uid")
			c.SetParamValues(tt.pathParams)
			err := s.setCanceledStatus(c)
			if err != nil {
				t.Errorf("getAllHotels() error = %v", err)
			}

			code := w.Result().StatusCode
			if code != tt.expectedHTTPStatus {
				t.Errorf("getAllHotels() http-code expected %d, but got %d", tt.expectedHTTPStatus, code)
			}
		})
	}
}
