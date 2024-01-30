package db

import "context"

// CreateUserTxParams contains the input parameters of the creat user transaction
type CreateUserTxParams struct {
	CreateUserParams
	AfterCreate func(user User) error
}

// CreateUserTxResult contains the result of the creat user transaction
type CreateUserTxResult struct {
	User User
}

// CreateUserTx creates a new user and send a verification email
// It creates a transfer record, add account entries and update account balance within a single database transaction
func (store *SQLStore) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error) {
	var result CreateUserTxResult

	err := store.execTx(ctx, func(queries *Queries) error {
		var err error

		result.User, err = queries.CreateUser(ctx, arg.CreateUserParams)

		if err != nil {
			return err
		}

		return arg.AfterCreate(result.User)
	})

	return result, err
}
