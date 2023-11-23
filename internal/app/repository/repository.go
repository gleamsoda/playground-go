package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/samber/do"

	"playground/internal/app"
)

type (
	Repository struct {
		exec Executor
		txn  app.Transaction
	}
	DB interface {
		BeginTx(context.Context, *sql.TxOptions) (*sql.Tx, error)
		Executor
	}
	Executor interface {
		ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
		PrepareContext(context.Context, string) (*sql.Stmt, error)
		QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
		QueryRowContext(context.Context, string, ...interface{}) *sql.Row
	}
)

func NewRepository(i *do.Injector) (app.Repository, error) {
	db := do.MustInvoke[*sql.DB](i)
	return &Repository{
		exec: db,
		txn:  NewTransaction(db),
	}, nil
}

func (r *Repository) Account() app.AccountRepository {
	return NewAccount(r.exec)
}

func (r *Repository) Transfer() app.TransferRepository {
	return NewTransfer(r.exec)
}

func (r *Repository) User() app.UserRepository {
	return NewUser(r.exec)
}

func (r *Repository) Transaction() app.Transaction {
	return r.txn
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
