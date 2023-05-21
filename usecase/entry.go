package usecase

import (
	"context"

	"github.com/gleamsoda/go-playground/domain"
)

type EntryUsecase struct {
	r domain.EntryRepository
}

var _ domain.EntryUsecase = (*EntryUsecase)(nil)

func NewEntryUsecase(r domain.EntryRepository) *EntryUsecase {
	return &EntryUsecase{
		r: r,
	}
}

func (u *EntryUsecase) Create(ctx context.Context, arg domain.CreateEntryParams) (*domain.Entry, error) {
	return u.r.Create(ctx, arg)
}

func (u *EntryUsecase) Get(ctx context.Context, id int64) (*domain.Entry, error) {
	return u.r.GetByID(ctx, id)
}
