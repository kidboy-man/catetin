package postgresql

import (
	"context"

	"github.com/ingunawandra/catetin/internal/repository"
	"gorm.io/gorm"
)

// gormDB is a thin wrapper around *gorm.DB that implements repository.DB
// so repositories depend on an abstraction instead of concrete GORM types.
type gormDB struct {
	db *gorm.DB
}

// NewDB wraps a *gorm.DB into repository.DB
func NewDB(db *gorm.DB) repository.DB {
	return &gormDB{db: db}
}

func (g *gormDB) WithContext(ctx context.Context) repository.DB {
	return &gormDB{db: g.db.WithContext(ctx)}
}

func (g *gormDB) Create(value interface{}) repository.Result {
	res := g.db.Create(value)
	return &gormResult{db: res}
}

func (g *gormDB) Where(query interface{}, args ...interface{}) repository.DB {
	return &gormDB{db: g.db.Where(query, args...)}
}

func (g *gormDB) First(dest interface{}) repository.Result {
	res := g.db.First(dest)
	return &gormResult{db: res}
}

func (g *gormDB) Limit(limit int) repository.DB {
	return &gormDB{db: g.db.Limit(limit)}
}

func (g *gormDB) Offset(offset int) repository.DB {
	return &gormDB{db: g.db.Offset(offset)}
}

func (g *gormDB) Order(value interface{}) repository.DB {
	return &gormDB{db: g.db.Order(value)}
}

func (g *gormDB) Find(dest interface{}) repository.Result {
	res := g.db.Find(dest)
	return &gormResult{db: res}
}

func (g *gormDB) Model(value interface{}) repository.DB {
	return &gormDB{db: g.db.Model(value)}
}

func (g *gormDB) Select(query interface{}) repository.DB {
	return &gormDB{db: g.db.Select(query)}
}

func (g *gormDB) Scan(dest interface{}) repository.Result {
	res := g.db.Scan(dest)
	return &gormResult{db: res}
}

func (g *gormDB) Updates(values interface{}) repository.Result {
	res := g.db.Updates(values)
	return &gormResult{db: res}
}

func (g *gormDB) Delete(value interface{}, conds ...interface{}) repository.Result {
	res := g.db.Delete(value, conds...)
	return &gormResult{db: res}
}

func (g *gormDB) Transaction(fn func(tx repository.DB) error) error {
	return g.db.Transaction(func(tx *gorm.DB) error {
		return fn(&gormDB{db: tx})
	})
}

func (g *gormDB) Begin() (repository.DB, error) {
	tx := g.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &gormDB{db: tx}, nil
}

func (g *gormDB) Commit() error {
	if err := g.db.Commit().Error; err != nil {
		return err
	}
	return nil
}

func (g *gormDB) Rollback() error {
	if err := g.db.Rollback().Error; err != nil {
		return err
	}
	return nil
}

// gormResult wraps *gorm.DB returned by query methods so we can expose
// the minimal Result interface to calling code.
type gormResult struct {
	db *gorm.DB
}

func (r *gormResult) Error() error        { return r.db.Error }
func (r *gormResult) RowsAffected() int64 { return r.db.RowsAffected }
