package persistence

import (
	"context"
	"errors"
)

// ErrTransactionClosed is returned by all methods on Transaction once the
// transaction is committed or rolled-back.
var ErrTransactionClosed = errors.New("transaction already committed or rolled-back")

// Transaction exposes persistence operations that can be performed atomically.
type Transaction interface {
	// AggregateTransaction
	// ProcessTransaction
	// QueueTransaction
	// eventstore.Transaction

	// Commit applies the changes from the transaction.
	Commit(ctx context.Context) error

	// Rollback aborts the transaction.
	Rollback() error
}
