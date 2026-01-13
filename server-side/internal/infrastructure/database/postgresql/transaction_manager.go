package postgresql

import (
	"context"
	"fmt"

	"github.com/ingunawandra/catetin/internal/repository"
	"gorm.io/gorm"
)

// transactionManager implements the TransactionManager interface using the
// repository.DB abstraction (backed by a GORM wrapper in this package).
type transactionManager struct {
	db repository.DB
}

// NewTransactionManager creates a new transaction manager. It accepts the
// concrete *gorm.DB and wraps it so the rest of the code can depend on
// repository.DB instead of *gorm.DB.
func NewTransactionManager(db *gorm.DB) repository.TransactionManager {
	return &transactionManager{db: NewDB(db)}
}

// WithTransaction executes a function within a database transaction
func (tm *transactionManager) WithTransaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	// If already in a transaction, just execute the function
	if tm.IsInTransaction(ctx) {
		return fn(ctx)
	}

	return tm.db.Transaction(func(tx repository.DB) error {
		// Create new context with transaction
		txCtx := repository.SetTransactionInContext(ctx, tx)
		return fn(txCtx)
	})
}

// BeginTransaction starts a new transaction
func (tm *transactionManager) BeginTransaction(ctx context.Context) (context.Context, error) {
	if tm.IsInTransaction(ctx) {
		return ctx, nil
	}

	tx, err := tm.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	txCtx := repository.SetTransactionInContext(ctx, tx)
	return txCtx, nil
}

// CommitTransaction commits the active transaction
func (tm *transactionManager) CommitTransaction(ctx context.Context) error {
	tx := repository.GetTransactionFromContext(ctx)
	if tx == nil {
		return fmt.Errorf("no active transaction to commit")
	}

	dbTx, ok := tx.(repository.DB)
	if !ok {
		return fmt.Errorf("invalid transaction type in context")
	}

	if err := dbTx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// RollbackTransaction rolls back the active transaction
func (tm *transactionManager) RollbackTransaction(ctx context.Context) error {
	tx := repository.GetTransactionFromContext(ctx)
	if tx == nil {
		return fmt.Errorf("no active transaction to rollback")
	}

	dbTx, ok := tx.(repository.DB)
	if !ok {
		return fmt.Errorf("invalid transaction type in context")
	}

	if err := dbTx.Rollback(); err != nil {
		return fmt.Errorf("failed to rollback transaction: %w", err)
	}

	return nil
}

// IsInTransaction checks if there's an active transaction in the context
func (tm *transactionManager) IsInTransaction(ctx context.Context) bool {
	return repository.GetTransactionFromContext(ctx) != nil
}

// GetDB returns the appropriate database connection (transaction or regular)
// This is a helper for repositories to use
func GetDB(ctx context.Context, db repository.DB) repository.DB {
	tx := repository.GetTransactionFromContext(ctx)
	if tx != nil {
		if dbTx, ok := tx.(repository.DB); ok {
			return dbTx
		}
	}
	return db.WithContext(ctx)
}
