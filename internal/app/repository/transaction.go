package repository

import (
	"context"
	"fmt"

	"playground/internal/app"
)

type (
	Transaction struct{ db TxRunner }
)

func NewTransaction(db TxRunner) app.Transaction {
	return &Transaction{db: db}
}

var _ app.Transaction = (*Transaction)(nil)

func (t *Transaction) Run(ctx context.Context, fn app.TransactionFunc) error {
	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	txm := &Manager{e: tx}
	if err := fn(ctx, txm); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}
