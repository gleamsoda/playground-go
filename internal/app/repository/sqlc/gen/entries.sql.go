// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: entries.sql

package gen

import (
	"context"
)

const createEntry = `-- name: CreateEntry :execlastid
INSERT INTO entries (
  account_id,
  amount
) VALUES (
  ?, ?
)
`

type CreateEntryParams struct {
	AccountID int64 `json:"account_id"`
	Amount    int64 `json:"amount"`
}

func (q *Queries) CreateEntry(ctx context.Context, arg *CreateEntryParams) (int64, error) {
	result, err := q.db.ExecContext(ctx, createEntry, arg.AccountID, arg.Amount)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

const getEntry = `-- name: GetEntry :one
SELECT id, account_id, amount, created_at FROM entries
WHERE id = ? LIMIT 1
`

func (q *Queries) GetEntry(ctx context.Context, id int64) (*Entry, error) {
	row := q.db.QueryRowContext(ctx, getEntry, id)
	var i Entry
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return &i, err
}

const listEntries = `-- name: ListEntries :many
SELECT id, account_id, amount, created_at FROM entries
WHERE account_id = ?
ORDER BY id
LIMIT ?
OFFSET ?
`

type ListEntriesParams struct {
	AccountID int64 `json:"account_id"`
	Limit     int32 `json:"limit"`
	Offset    int32 `json:"offset"`
}

func (q *Queries) ListEntries(ctx context.Context, arg *ListEntriesParams) ([]*Entry, error) {
	rows, err := q.db.QueryContext(ctx, listEntries, arg.AccountID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*Entry{}
	for rows.Next() {
		var i Entry
		if err := rows.Scan(
			&i.ID,
			&i.AccountID,
			&i.Amount,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
