package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

// TxFunc тип функции, которая будет выполняться внутри транзакции
type TxFunc func(ctx context.Context, tx *sql.Tx) error

// WithWriteTransaction обертка для управления транзакцией
func WithWriteTransaction(ctx context.Context, db *sqlx.DB, fn TxFunc) error {
	return withTx(ctx, db, fn, false)
}

func withTx(ctx context.Context, db *sqlx.DB, fn TxFunc, isReadOnly bool) error {
	tx, err := db.BeginTx(ctx, &sql.TxOptions{
		ReadOnly: isReadOnly,
	})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	committed := false

	defer func() {
		if !committed {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("failed to rollback transaction: %v", rollbackErr)
			}
		}
	}()

	if err := fn(ctx, tx); err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	committed = true
	return nil
}

// WithNoTransaction - обертка с R/O транзакцией
func WithNoTransaction(ctx context.Context, db *sqlx.DB, fn TxFunc) error {
	return withTx(ctx, db, fn, true)
}
