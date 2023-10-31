package app

import (
	"time"
)

type Entry struct {
	ID        int64 `json:"id"`
	AccountID int64 `json:"account_id"`
	// can be negative or positive
	Amount    int64     `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

func NewEntry(accountID, amount int64) *Entry {
	return &Entry{
		AccountID: accountID,
		Amount:    amount,
	}
}
