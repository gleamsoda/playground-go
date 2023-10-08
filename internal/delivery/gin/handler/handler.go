package handler

import (
	"github.com/samber/do"

	"playground/internal/wallet"
)

type Handler struct {
	w wallet.Usecase
}

func NewHandler(i *do.Injector) (*Handler, error) {
	w := do.MustInvoke[wallet.Usecase](i)
	return &Handler{w: w}, nil
}
