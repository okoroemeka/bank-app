package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
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

	user, err := testStore.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.True(t, user.PasswordChangedAt.Time.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomAccount(t)
}

func TestGetUser(t *testing.T) {
	randomUser := createRandomUser(t)
	user, err := testStore.GetUser(context.Background(), randomUser.Username)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, randomUser.Email, user.Email)
	require.Equal(t, randomUser.FullName, user.FullName)
	require.Equal(t, randomUser.HashedPassword, user.HashedPassword)
	require.Equal(t, randomUser.Username, user.Username)
	require.WithinDuration(t, randomUser.PasswordChangedAt.Time, user.PasswordChangedAt.Time, time.Second)
	require.WithinDuration(t, randomUser.CreatedAt.Time, user.CreatedAt.Time, time.Second)
}

func TestUpdateUser(t *testing.T) {
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

//func TestDeleteAccount(t *testing.T) {
//	account1 := createRandomAccount(t)
//
//	err := testStore.DeleteAccount(context.Background(), account1.ID)
//	require.NoError(t, err)
//
//	account2, err := testStore.GetAccount(context.Background(), account1.ID)
//	require.Error(t, err)
//	require.EqualError(t, err, ErrorNoRecordFound.Error())
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
//	accounts, err := testStore.ListAccounts(context.Background(), params)
//
//	require.NoError(t, err)
//	require.Len(t, accounts, 5)
//
//	for _, account := range accounts {
//		require.NotEmpty(t, account)
//	}
//}

func TestUpdateOnlyUserFullname(t *testing.T) {
	oldUser := createRandomUser(t)

	newFullName := util.RandomOwner()
	updatedUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		FullName: pgtype.Text{String: newFullName, Valid: true},
	})

	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)
	require.NotEqual(t, oldUser.FullName, updatedUser.FullName)
	require.Equal(t, newFullName, updatedUser.FullName)
	require.Equal(t, oldUser.Email, updatedUser.Email)
}
func TestUpdateOnlyUserEmail(t *testing.T) {
	oldUser := createRandomUser(t)

	newEmail := util.RandomEmail()
	updatedUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		Email:    pgtype.Text{String: newEmail, Valid: true},
	})

	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)
	require.NotEqual(t, oldUser.Email, updatedUser.Email)
	require.Equal(t, newEmail, updatedUser.Email)
	require.Equal(t, oldUser.FullName, updatedUser.FullName)
}
func TestUpdateOnlyUserPassword(t *testing.T) {
	oldUser := createRandomUser(t)

	newPassword, err := util.HashPassword(util.RandomString(10))

	require.NoError(t, err)

	updatedUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		Username:       oldUser.Username,
		HashedPassword: pgtype.Text{String: newPassword, Valid: true},
	})

	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)
	require.NotEqual(t, oldUser.HashedPassword, updatedUser.HashedPassword)
	require.Equal(t, newPassword, updatedUser.HashedPassword)
	require.Equal(t, oldUser.FullName, updatedUser.FullName)
	require.Equal(t, oldUser.Email, updatedUser.Email)
}

func TestUpdateUserFields(t *testing.T) {
	oldUser := createRandomUser(t)

	newFullName := util.RandomOwner()
	newEmail := util.RandomEmail()
	newHashedPassword, err := util.HashPassword(util.RandomString(10))

	require.NoError(t, err)

	updatedUser, err := testStore.UpdateUser(context.Background(), UpdateUserParams{
		Username:       oldUser.Username,
		FullName:       pgtype.Text{String: newFullName, Valid: true},
		Email:          pgtype.Text{String: newEmail, Valid: true},
		HashedPassword: pgtype.Text{String: newHashedPassword, Valid: true},
	})

	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)
	require.NotEqual(t, oldUser.FullName, updatedUser.FullName)
	require.NotEqual(t, oldUser.Email, updatedUser.Email)
	require.NotEqual(t, oldUser.HashedPassword, updatedUser.HashedPassword)
	require.Equal(t, newFullName, updatedUser.FullName)
	require.Equal(t, newEmail, updatedUser.Email)
	require.Equal(t, newHashedPassword, updatedUser.HashedPassword)
}
