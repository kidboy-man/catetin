package postgresql

import (
	"context"
	"fmt"

	"github.com/ingunawandra/catetin/internal/repository"
	"gorm.io/gorm"
)

// transactionManager implements the TransactionManager interface using GORM
type transactionManager struct {
	db *gorm.DB
}

// NewTransactionManager creates a new GORM-based transaction manager
func NewTransactionManager(db *gorm.DB) repository.TransactionManager {
	return &transactionManager{db: db}
}

// WithTransaction executes a function within a database transaction
func (tm *transactionManager) WithTransaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	// Check if already in a transaction
	if tm.IsInTransaction(ctx) {
		// Already in a transaction, just execute the function
		return fn(ctx)
	}

	// Start a new transaction
	return tm.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create new context with transaction
		txCtx := repository.SetTransactionInContext(ctx, tx)

		// Execute the function
		return fn(txCtx)
	})
}

// BeginTransaction starts a new transaction
func (tm *transactionManager) BeginTransaction(ctx context.Context) (context.Context, error) {
	// Check if already in a transaction
	if tm.IsInTransaction(ctx) {
		return ctx, nil
	}

	// Begin transaction
	tx := tm.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	// Store transaction in context
	txCtx := repository.SetTransactionInContext(ctx, tx)
	return txCtx, nil
}

// CommitTransaction commits the active transaction
func (tm *transactionManager) CommitTransaction(ctx context.Context) error {
	tx := repository.GetTransactionFromContext(ctx)
	if tx == nil {
		return fmt.Errorf("no active transaction to commit")
	}

	gormTx, ok := tx.(*gorm.DB)
	if !ok {
		return fmt.Errorf("invalid transaction type in context")
	}

	if err := gormTx.Commit().Error; err != nil {
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

	gormTx, ok := tx.(*gorm.DB)
	if !ok {
		return fmt.Errorf("invalid transaction type in context")
	}

	if err := gormTx.Rollback().Error; err != nil {
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
func GetDB(ctx context.Context, db *gorm.DB) *gorm.DB {
	tx := repository.GetTransactionFromContext(ctx)
	if tx != nil {
		if gormTx, ok := tx.(*gorm.DB); ok {
			return gormTx
		}
	}
	return db.WithContext(ctx)
}
