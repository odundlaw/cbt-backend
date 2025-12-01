// Package validation for handling valiation
package validation

import (
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
	response "github.com/odundlaw/cbt-backend/internal/json"
)

var Validate *validator.Validate

func init() {
	v := validator.New()

	// Use JSON field names instead of struct field names
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := fld.Tag.Get("json")
		if name == "-" {
			return ""
		}
		return name
	})

	Validate = v
}

func FormatValidationErrors(err error) []response.FieldError {
	errs := []response.FieldError{}

	for _, fe := range err.(validator.ValidationErrors) {
		msg := ""

		switch fe.Tag() {
		case "required":
			msg = fmt.Sprintf("%s is required", fe.Field())
		case "email":
			msg = "Invalid email format"
		case "min":
			msg = fmt.Sprintf("%s must be at least %s characters", fe.Field(), fe.Param())
		case "max":
			msg = fmt.Sprintf("%s must not be more than %s characters", fe.Field(), fe.Param())
		default:
			msg = fmt.Sprintf("%s is not valid", fe.Field())
		}

		errs = append(errs, response.FieldError{
			Field:   fe.Field(),
			Message: msg,
		})
	}

	return errs
}
