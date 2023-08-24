package app

import "context"

type TransactionManager interface {
	Do(ctx context.Context, fn func(context.Context) error) error
}

type PassThroughTxManager struct{}

var _ TransactionManager = (*PassThroughTxManager)(nil)

func (m *PassThroughTxManager) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}
