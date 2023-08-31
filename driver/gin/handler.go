package gin

import "playground/app"

type handler struct {
	u app.Usecase
}

func NewHandler(u app.Usecase) *handler {
	return &handler{
		u: u,
	}
}
