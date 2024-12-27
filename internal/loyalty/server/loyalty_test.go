package server

import (
	"encoding/json"
	"github.com/AskaryanKarine/bmstu-ds-4/pkg/models"
	"github.com/AskaryanKarine/bmstu-ds-4/pkg/validation"
	"github.com/go-playground/validator/v10"
	"github.com/gojuno/minimock/v3"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestServer_getLoyaltyByUser(t *testing.T) {
	mc := minimock.NewController(t)
	e := echo.New()
	e.Validator = validation.MustRegisterCustomValidator(validator.New())
	type fields struct {
		echo *echo.Echo
		lr   loyaltyRepository
	}
	regularDiscount := 5
	regularLoyalty := models.LoyaltyInfoResponse{
		Status:           models.BRONZE,
		Discount:         5,
		ReservationCount: &regularDiscount,
	}

	tests := []struct {
		name               string
		fields             fields
		expectedHTTPStatus int
		result             models.LoyaltyInfoResponse
		userHeader         string
	}{
		{
			name: "http-200",
			fields: fields{
				echo: e,
				lr:   NewLoyaltyRepositoryMock(mc).GetByUserMock.Return(regularLoyalty, nil),
			},
			expectedHTTPStatus: http.StatusOK,
			result:             regularLoyalty,
			userHeader:         "Test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				echo: tt.fields.echo,
				lr:   tt.fields.lr,
			}
			r := httptest.NewRequest(http.MethodGet, "/test", nil)
			w := httptest.NewRecorder()
			c := s.echo.NewContext(r, w)
			c.Set("username", tt.userHeader)
			err := s.getLoyaltyByUser(c)
			if err != nil {
				t.Errorf("getLoyaltyByUser() error = %v", err)
			}

			code := w.Result().StatusCode
			if code != tt.expectedHTTPStatus {
				t.Errorf("getLoyaltyByUser() http-code expected %d, but got %d", tt.expectedHTTPStatus, code)
			}

			body, err := io.ReadAll(w.Result().Body)
			if err != nil {
				t.Errorf("ReadAll error")
			}
			var res models.LoyaltyInfoResponse
			err = json.Unmarshal(body, &res)
			if err != nil {
				t.Errorf("json unmarshal error")
			}

			if tt.expectedHTTPStatus == http.StatusOK {
				if !reflect.DeepEqual(res, tt.result) {
					t.Errorf("getLoyaltyByUser() expected %v, but got %v", tt.result, res)
				}
			}
		})
	}
}
