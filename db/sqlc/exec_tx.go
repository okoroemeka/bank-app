package db

import (
	"context"
	"fmt"
)

// execTx executes a function or performs a unit of work within a database transaction(unit of work)
func (store *SQLStore) execTx(ctx context.Context, fn func(queries *Queries) error) error {
	tx, err := store.connPool.Begin(ctx)

	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx error %v, rb err %v", err, rbErr)
		}
		return err
	}
	return tx.Commit(ctx)
}
