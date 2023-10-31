package resolver

import (
	"github.com/samber/do"

	"playground/internal/app"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	w app.Usecase
}

func NewResolver(i *do.Injector) (*Resolver, error) {
	w := do.MustInvoke[app.Usecase](i)
	return &Resolver{w: w}, nil
}
