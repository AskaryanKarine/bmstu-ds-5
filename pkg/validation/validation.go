package validation

import (
	"github.com/go-playground/validator/v10"
	"time"
)

type customValidator struct {
	validator *validator.Validate
}

func MustRegisterCustomValidator(v *validator.Validate) *customValidator {
	err := v.RegisterValidation("IsISO8601", validateISO8601)
	if err != nil {

	}

	return &customValidator{validator: v}
}

func (cv *customValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func validateISO8601(fl validator.FieldLevel) bool {
	value := fl.Field().Interface()
	if value == nil {
		return false
	}

	date, ok := value.(string)
	if !ok {
		return false
	}

	_, err := time.Parse("2006-01-02", date)

	return err == nil
}
