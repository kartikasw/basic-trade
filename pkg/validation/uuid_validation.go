package validation

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var ValidUUID validator.Func = func(fl validator.FieldLevel) bool {
	if value, ok := fl.Field().Interface().(string); ok {
		_, err := uuid.Parse(value)
		if err != nil {
			return false
		}

		return true
	}

	return false
}
