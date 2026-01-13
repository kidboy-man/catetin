package repository

import "context"

// DB is a minimal database abstraction used by repositories to avoid coupling
// to a concrete ORM implementation (e.g., GORM). It intentionally mirrors the
// methods used across repository implementations.
type DB interface {
	WithContext(ctx context.Context) DB
	Create(value interface{}) Result
	Where(query interface{}, args ...interface{}) DB
	First(dest interface{}) Result
	Limit(limit int) DB
	Offset(offset int) DB
	Order(value interface{}) DB
	Find(dest interface{}) Result
	Model(value interface{}) DB
	Select(query interface{}) DB
	Scan(dest interface{}) Result
	Updates(values interface{}) Result
	Delete(value interface{}, conds ...interface{}) Result

	// Transaction helpers
	Transaction(fn func(tx DB) error) error
	Begin() (DB, error)
	Commit() error
	Rollback() error
}

// Result abstracts the execution result (Error and RowsAffected) returned by
// query operations.
type Result interface {
	Error() error
	RowsAffected() int64
}
