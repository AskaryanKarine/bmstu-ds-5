package validation

import (
	"errors"
	"github.com/AskaryanKarine/bmstu-ds-4/pkg/models"
	"github.com/go-playground/validator/v10"
)

func ConvertToError(err error, comment string) error {
	var valErr validator.ValidationErrors
	if errors.As(err, &valErr) {
		var respErr models.ValidationErrorResponse
		respErr.Message = comment
		respErr.Errors = make([]models.ErrorDescription, 0, len(valErr))
		for i := range valErr {
			respErr.Errors = append(respErr.Errors, models.ErrorDescription{
				Field: valErr[i].Field(),
				Error: valErr[i].Error(),
			})
		}
		return respErr
	}
	return err
}
