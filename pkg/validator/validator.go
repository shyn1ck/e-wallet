package validator

import (
	apperrors "e-wallet/pkg/errors"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// Validate a struct
func Validate(s interface{}) error {
	err := validate.Struct(s)
	if err != nil {
		return apperrors.ErrValidationFailed
	}
	return nil
}
