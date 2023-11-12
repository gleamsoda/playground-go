package repository

import (
	"context"
	"testing"
	"time"

	"github.com/morikuni/failure"
	"github.com/stretchr/testify/require"

	"playground/internal/app"
	"playground/internal/pkg/apperr"
)

func createRandomAccount(t *testing.T) *app.Account {
	user := createRandomUser(t)

	arg := &app.Account{
		Owner:    user.Username,
		Balance:  app.RandomMoney(),
		Currency: app.RandomCurrency(),
	}

	account, err := rm.Account().Create(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := rm.Account().Get(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, account1.Balance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	newBalance := app.RandomMoney()
	account1.Balance = newBalance
	account2, err := rm.Account().Update(context.Background(), account1)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account1.Owner, account2.Owner)
	require.Equal(t, newBalance, account2.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	err := rm.Account().Delete(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err := rm.Account().Get(context.Background(), account1.ID)
	require.Error(t, err)
	code, _ := failure.CodeOf(err)
	require.Equal(t, apperr.NotFound, code)
	require.Empty(t, account2)
}

func TestListAccounts(t *testing.T) {
	var lastAccount *app.Account
	for i := 0; i < 10; i++ {
		lastAccount = createRandomAccount(t)
	}

	arg := &app.ListAccountsParams{
		Owner:  lastAccount.Owner,
		Limit:  5,
		Offset: 0,
	}

	accounts, err := rm.Account().List(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.Equal(t, lastAccount.Owner, account.Owner)
	}
}
