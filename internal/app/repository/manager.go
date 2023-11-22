package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/samber/do"

	"playground/internal/app"
)

type (
	Manager struct {
		r TxRunner
		e Executor
	}
	TxRunner interface {
		BeginTx(context.Context, *sql.TxOptions) (*sql.Tx, error)
	}
	Executor interface {
		ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
		PrepareContext(context.Context, string) (*sql.Stmt, error)
		QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
		QueryRowContext(context.Context, string, ...interface{}) *sql.Row
	}
)

func NewManager(i *do.Injector) (app.RepositoryManager, error) {
	db := do.MustInvoke[*sql.DB](i)
	return &Manager{
		r: db,
		e: db,
	}, nil
}

func (m *Manager) Account() app.AccountRepository {
	return NewAccount(m.e)
}

func (m *Manager) Transfer() app.TransferRepository {
	return NewTransfer(m.e)
}

func (m *Manager) User() app.UserRepository {
	return NewUser(m.e)
}

func (m *Manager) Transaction() app.Transaction {
	return NewTransaction(m.r)
}

func runTx(ctx context.Context, e Executor, fn func(context.Context, *sql.Tx) error) error {
	switch v := e.(type) {
	case *sql.Tx:
		return fn(ctx, v)
	case *sql.DB:
		tx, err := v.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		if err := fn(ctx, tx); err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
			}
			return err
		}
		return tx.Commit()
	default:
		panic("invalid db type")
	}
}
