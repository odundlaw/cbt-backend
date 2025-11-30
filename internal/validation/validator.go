// Package validation for handling valiation
package validation

import (
	"reflect"

	"github.com/go-playground/validator/v10"
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
