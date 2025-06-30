package storage

import (
	"botmediasaver/generated/sqlc"
	"botmediasaver/internal/client/pgxpool"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Storage struct {
	db pgxpool.DBTX
	*sqlc.Queries
}

type TxStorage struct {
	*Storage
	tx pgx.Tx
	// Add more tx fields if needed (e.g., redis, mongo, etc.)
}

func NewStorage(db pgxpool.DBTX) *Storage {
	return &Storage{
		db:      db,
		Queries: sqlc.New(db),
	}
}

// BeginTx starts a pseudo nested transaction.
func (s *Storage) BeginTx(ctx context.Context) (*TxStorage, error) {
	// Add more tx begin logics if needed (e.g., redis, mongo, etc.)
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	return &TxStorage{
		Storage: NewStorage(tx),
		tx:      tx,
	}, nil
}

// Commit commits the transaction if this is a real transaction or releases the savepoint if this is a pseudo nested
// transaction. Commit will return an error if the Tx is already closed, but is otherwise safe to call multiple times.
func (ts *TxStorage) Commit(ctx context.Context) error {
	// Add more tx commit logics if needed (e.g., redis, mongo, etc.)
	return ts.tx.Commit(ctx)
}

// Rollback rolls back the transaction. Rollback will return an error if the Tx is already closed, but is otherwise safe
// to call multiple times. Hence, a defer storage.Rollback() is safe (must safe) even if storage.Commit() will be
// called first in a non-error condition. Any other failure of a real transaction will result in the connection being closed.
func (ts *TxStorage) Rollback(ctx context.Context) error {
	// Add more tx rollback logics if needed (e.g., redis, mongo, etc.)
	return ts.tx.Rollback(ctx)
}
