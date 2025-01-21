package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	Querier
}

type DBStore struct {
	*Queries
	conn *pgxpool.Pool
}

func NewStore(connPool *pgxpool.Pool) Store {
	return &DBStore{
		Queries: New(connPool),
		conn:    connPool,
	}
}

func (store *DBStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("transaction error: %v, rollback error: %v", err, rbErr)
		}
		return fmt.Errorf("transaction error: %v", err)
	}

	return tx.Commit(ctx)
}
