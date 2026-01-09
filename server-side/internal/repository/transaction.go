package repository

import "context"

// TransactionManager defines the interface for managing database transactions
// This abstraction allows the service layer to use transactions without knowing
// the underlying database implementation (GORM, sqlx, etc.)
type TransactionManager interface {
	// WithTransaction executes the given function within a database transaction.
	// If the function returns an error, the transaction is rolled back.
	// Otherwise, the transaction is committed.
	//
	// Example usage in service layer:
	//   err := txManager.WithTransaction(ctx, func(txCtx context.Context) error {
	//       // All repository calls using txCtx will use the same transaction
	//       if err := userRepo.Create(txCtx, user); err != nil {
	//           return err // will trigger rollback
	//       }
	//       if err := userAuthRepo.Create(txCtx, userAuth); err != nil {
	//           return err // will trigger rollback
	//       }
	//       return nil // will commit transaction
	//   })
	WithTransaction(ctx context.Context, fn func(txCtx context.Context) error) error

	// BeginTransaction starts a new transaction and returns a context with the transaction.
	// This is useful when you need more control over the transaction lifecycle.
	// Remember to call CommitTransaction or RollbackTransaction.
	BeginTransaction(ctx context.Context) (context.Context, error)

	// CommitTransaction commits the transaction associated with the context.
	CommitTransaction(ctx context.Context) error

	// RollbackTransaction rolls back the transaction associated with the context.
	RollbackTransaction(ctx context.Context) error

	// IsInTransaction checks if the context has an active transaction
	IsInTransaction(ctx context.Context) bool
}

// TransactionKey is the context key for storing transaction information
type transactionKey struct{}

// TxKey is the actual key instance used for context
var TxKey = transactionKey{}

// GetTransactionFromContext retrieves the transaction from context
// Returns nil if no transaction is active
func GetTransactionFromContext(ctx context.Context) interface{} {
	return ctx.Value(TxKey)
}

// SetTransactionInContext stores the transaction in context
func SetTransactionInContext(ctx context.Context, tx interface{}) context.Context {
	return context.WithValue(ctx, TxKey, tx)
}
