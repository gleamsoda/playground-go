package internal

import (
	"github.com/go-playground/validator/v10"

	"github.com/gleamsoda/go-playground/domain"
)

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if c, ok := fieldLevel.Field().Interface().(string); ok {
		return domain.IsSupportedCurrency(c)
	}
	return false
}
