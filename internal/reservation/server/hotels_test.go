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

var (
	regularHotelPagination = models.PaginationResponse{
		Page:          1,
		PageSize:      1,
		TotalElements: 0,
		Items:         nil,
	}
)

func TestServer_getAllHotels(t *testing.T) {
	mc := minimock.NewController(t)
	e := echo.New()
	e.Validator = validation.MustRegisterCustomValidator(validator.New())

	type fields struct {
		echo *echo.Echo
		hs   hotelStorage
		rs   reservationStorage
	}
	tests := []struct {
		name               string
		fields             fields
		expectedHTTPStatus int
		result             models.PaginationResponse
		queryParams        string
	}{
		{
			name: "http-200: success response",
			fields: fields{
				echo: e,
				hs:   NewHotelStorageMock(mc).GetAllHotelsMock.Return(nil, 0, nil),
				rs:   nil,
			},
			expectedHTTPStatus: http.StatusOK,
			result:             regularHotelPagination,
			queryParams:        "?page=1&size=1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				echo: tt.fields.echo,
				hs:   tt.fields.hs,
				rs:   tt.fields.rs,
			}
			r := httptest.NewRequest(http.MethodGet, "/test"+tt.queryParams, nil)
			w := httptest.NewRecorder()
			c := s.echo.NewContext(r, w)
			err := s.getAllHotels(c)
			if err != nil {
				t.Errorf("getAllHotels() error = %v", err)
			}

			code := w.Result().StatusCode
			if code != tt.expectedHTTPStatus {
				t.Errorf("getAllHotels() http-code expected %d, but got %d", tt.expectedHTTPStatus, code)
			}

			body, err := io.ReadAll(w.Result().Body)
			if err != nil {
				t.Errorf("ReadAll error")
			}
			var res models.PaginationResponse
			err = json.Unmarshal(body, &res)
			if err != nil {
				t.Errorf("json unmarshal error")
			}

			if tt.expectedHTTPStatus == http.StatusOK {
				if !reflect.DeepEqual(res, tt.result) {
					t.Errorf("getPersonByID() expected %v, but got %v", tt.result, res)
				}
			}
		})
	}
}
