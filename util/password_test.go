package util

import (
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestHashPassword(t *testing.T) {
	randomPass := RandomString(10)

	hashedPass, err := HashPassword(randomPass)

	require.NoError(t, err)
	require.NotEmpty(t, hashedPass)

	err = ComparePassword(randomPass, hashedPass)
	require.NoError(t, err)

	wrongPass := RandomString(9)

	err = ComparePassword(wrongPass, hashedPass)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	hashedPass1, err := HashPassword(randomPass)

	require.NoError(t, err)
	require.NotEmpty(t, hashedPass1)
	require.NotEqual(t, hashedPass, hashedPass1)
}
