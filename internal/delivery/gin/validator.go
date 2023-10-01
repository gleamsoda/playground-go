package gin

import (
	"github.com/go-playground/validator/v10"

	"playground/internal/wallet"
)

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if c, ok := fieldLevel.Field().Interface().(string); ok {
		return wallet.IsSupportedCurrency(c)
	}
	return false
}
