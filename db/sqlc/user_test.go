package db

import (
	"context"
	"github.com/okoroemeka/simple_bank/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createRandomUser(t *testing.T) User {
	hashedPassed, err := util.HashPassword(util.RandomString(6))

	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
		HashedPassword: hashedPassed,
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomAccount(t)
}

func TestGetUser(t *testing.T) {
	randomUser := createRandomUser(t)
	user, err := testQueries.GetUser(context.Background(), randomUser.Username)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, randomUser.Email, user.Email)
	require.Equal(t, randomUser.FullName, user.FullName)
	require.Equal(t, randomUser.HashedPassword, user.HashedPassword)
	require.Equal(t, randomUser.Username, user.Username)
	require.WithinDuration(t, randomUser.PasswordChangedAt, user.PasswordChangedAt, time.Second)
	require.WithinDuration(t, randomUser.CreatedAt, user.CreatedAt, time.Second)
}

//func TestUpdateUser(t *testing.T) {
//	account1 := createRandomAccount(t)
//	arg := UpdateAccountParams{
//		ID:      account1.ID,
//		Balance: util.RandomMoney(),
//	}
//	account2, err := testQueries.UpdateAccount(context.Background(), arg)
//	require.NoError(t, err)
//	require.NotEmpty(t, account2)
//	require.NotEqual(t, account2.Balance, account1.Balance)
//	require.Equal(t, account1.ID, account2.ID)
//	require.Equal(t, account2.Owner, account1.Owner)
//}

//func TestDeleteAccount(t *testing.T) {
//	account1 := createRandomAccount(t)
//
//	err := testQueries.DeleteAccount(context.Background(), account1.ID)
//	require.NoError(t, err)
//
//	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
//	require.Error(t, err)
//	require.EqualError(t, err, sql.ErrNoRows.Error())
//	require.Empty(t, account2)
//}
//
//func TestListAccount(t *testing.T) {
//	for i := 0; i < 10; i++ {
//		createRandomAccount(t)
//	}
//
//	params := ListAccountsParams{
//		Limit:  5,
//		Offset: 5,
//	}
//	accounts, err := testQueries.ListAccounts(context.Background(), params)
//
//	require.NoError(t, err)
//	require.Len(t, accounts, 5)
//
//	for _, account := range accounts {
//		require.NotEmpty(t, account)
//	}
//}
