package db

import (
	"context"
	customerror "github.com/okoroemeka/simple_bank/custom-error"
	"github.com/okoroemeka/simple_bank/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testStore.CreateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Currency, account.Currency)
	require.Equal(t, arg.Balance, account.Balance)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}
func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testStore.GetAccount(context.Background(), account1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account2.Owner, account1.Owner)
	require.Equal(t, account2.Currency, account1.Currency)
	require.Equal(t, account2.Balance, account1.Balance)
	require.Equal(t, account2.ID, account1.ID)
}

func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: util.RandomMoney(),
	}
	account2, err := testStore.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.NotEqual(t, account2.Balance, account1.Balance)
	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account2.Owner, account1.Owner)
}

func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	err := testStore.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	account2, err := testStore.GetAccount(context.Background(), account1.ID)
	require.Error(t, err)
	require.EqualError(t, err, customerror.ErrorNoRecordFound.Error())
	require.Empty(t, account2)
}

func TestListAccount(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	params := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}
	accounts, err := testStore.ListAccounts(context.Background(), params)

	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
