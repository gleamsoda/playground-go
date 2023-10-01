package repository

import (
	"context"
	"database/sql"
	"fmt"

	"playground/internal/wallet"
	"playground/internal/wallet/repository/sqlc/gen"
)

type Repository struct {
	db *sql.DB
	q  *gen.Queries
}

func NewRepository(db *sql.DB) wallet.Repository {
	return &Repository{
		db: db,
		q:  gen.New(db),
	}
}

func (r *Repository) tx(ctx context.Context, fn func(context.Context, *gen.Queries) error) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := gen.New(tx)
	if err := fn(ctx, q); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}
