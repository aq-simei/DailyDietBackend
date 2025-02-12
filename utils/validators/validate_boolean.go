package validators

import (
	"github.com/go-playground/validator/v10"
)

// ValidateBoolean ensures that the boolean field is explicitly set
func ValidateBoolean(fl validator.FieldLevel) bool {
	_, ok := fl.Field().Interface().(bool)
	return ok
}
