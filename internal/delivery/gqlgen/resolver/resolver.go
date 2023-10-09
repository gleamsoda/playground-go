package resolver

import (
	"github.com/samber/do"

	"playground/internal/wallet"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	w wallet.Usecase
}

func NewResolver(i *do.Injector) (*Resolver, error) {
	w := do.MustInvoke[wallet.Usecase](i)
	return &Resolver{w: w}, nil
}
