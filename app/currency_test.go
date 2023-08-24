package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsSupportedCurrency(t *testing.T) {
	supportedCurrencies := []string{"USD", "EUR", "CAD"}
	unsupportedCurrencies := []string{"JPY", "GBP", "AUD"}

	for _, currency := range supportedCurrencies {
		assert.True(t, IsSupportedCurrency(currency), "Expected IsSupportedCurrency(%s) to be true, but it was false", currency)
	}
	for _, currency := range unsupportedCurrencies {
		assert.False(t, IsSupportedCurrency(currency), "Expected IsSupportedCurrency(%s) to be false, but it was true", currency)
	}
}
