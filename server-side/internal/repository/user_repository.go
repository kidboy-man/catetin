package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ingunawandra/catetin/internal/domain"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *domain.User) error

	// FindByID finds a user by ID
	FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error)

	// FindByPhoneNumber finds a user by phone number
	FindByPhoneNumber(ctx context.Context, phoneNumber string) (*domain.User, error)

	// Update updates an existing user
	Update(ctx context.Context, user *domain.User) error

	// Delete soft deletes a user
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves all users with pagination
	List(ctx context.Context, limit, offset int) ([]*domain.User, error)
}
