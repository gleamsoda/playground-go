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

type ListEntriesParams struct {
	RequestUserID int64 `json:"-"`
	AccountID     int64 `json:"account_id"`
	Limit         int32 `json:"limit" form:"limit"`
	Offset        int32 `json:"offset" form:"offset"`
}
