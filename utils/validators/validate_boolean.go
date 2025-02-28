package validators

import (
	"reflect"

	"github.com/go-playground/validator/v10"
)

// ValidateBoolean ensures that the boolean field is explicitly set
func ValidateBoolean(fl validator.FieldLevel) bool {
	field := fl.Field()
	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			return false
		}
		field = field.Elem()
	}
	_, ok := field.Interface().(bool)
	return ok
}
