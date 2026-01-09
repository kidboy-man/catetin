package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/ingunawandra/catetin/internal/domain"
)

// MoneyFlowRepository defines the interface for money flow data access
type MoneyFlowRepository interface {
	// Create creates a new money flow
	Create(ctx context.Context, moneyFlow *domain.MoneyFlow) error

	// FindByID finds a money flow by ID
	FindByID(ctx context.Context, id uuid.UUID) (*domain.MoneyFlow, error)

	// FindByUserID finds all money flows for a specific user
	FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*domain.MoneyFlow, error)

	// FindByUserIDAndDateRange finds money flows for a user within a date range
	FindByUserIDAndDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]*domain.MoneyFlow, error)

	// Update updates an existing money flow
	Update(ctx context.Context, moneyFlow *domain.MoneyFlow) error

	// Delete soft deletes a money flow
	Delete(ctx context.Context, id uuid.UUID) error

	// GetTotalByUserID calculates total expenses for a user
	GetTotalByUserID(ctx context.Context, userID uuid.UUID) (float64, error)

	// GetTotalByUserIDAndCategory calculates total expenses by category
	GetTotalByUserIDAndCategory(ctx context.Context, userID uuid.UUID, category string) (float64, error)
}
