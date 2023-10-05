package handler

import "playground/internal/wallet"

type Handler struct {
	w wallet.Usecase
}

func NewHandler(w wallet.Usecase) *Handler {
	return &Handler{w: w}
}
