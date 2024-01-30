package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
)

type VerifyEmailTxParams struct {
	ID   int64
	Code string
}

type VerifyEmailTxResult struct {
	User        User
	VerifyEmail VerifyEmail
}

func (store *SQLStore) VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (result VerifyEmailTxResult, err error) {

	err = store.execTx(ctx, func(queries *Queries) error {

		verifyEmail, err := queries.UpdateVerifyEmailIsUsedField(ctx, UpdateVerifyEmailIsUsedFieldParams{
			IsUsed:     true,
			ID:         arg.ID,
			SecretCode: arg.Code,
		})

		if err != nil {
			return err
		}

		result.VerifyEmail = verifyEmail

		user, err := queries.UpdateUser(ctx, UpdateUserParams{
			Username:        verifyEmail.Username,
			IsEmailVerified: pgtype.Bool{Bool: true, Valid: true},
		})

		if err != nil {
			return err
		}

		result.User = user

		return nil
	})

	return result, err
}
