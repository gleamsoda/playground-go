package server

import (
	"github.com/go-playground/validator/v10"

	"playground/app"
)

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if c, ok := fieldLevel.Field().Interface().(string); ok {
		return app.IsSupportedCurrency(c)
	}
	return false
}
