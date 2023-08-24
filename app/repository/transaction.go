package repository

import (
	"database/sql"

	"golang.org/x/net/context"

	"playground/app"
)

type TransactionManager struct {
	db *sql.DB
}

var _ app.TransactionManager = (*TransactionManager)(nil)

const TransactionKey = "txKey"

func NewTransactionManager(db *sql.DB) app.TransactionManager {
	return &TransactionManager{
		db: db,
	}
}

func (r *TransactionManager) Do(ctx context.Context, fn func(context.Context) error) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	newCtx := context.WithValue(ctx, TransactionKey, tx)

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		}
	}()

	if err := fn(newCtx); err != nil {
		return err
	}

	return tx.Commit()
}
