package handler

import (
	"github.com/samber/do"

	"playground/internal/app"
)

type Handler struct {
	w app.Usecase
}

func NewHandler(i *do.Injector) (*Handler, error) {
	w := do.MustInvoke[app.Usecase](i)
	return &Handler{w: w}, nil
}
