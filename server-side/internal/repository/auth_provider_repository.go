package repository

import (
	"context"

	"github.com/google/uuid"
)

// AuthProvider represents an authentication provider entity
type AuthProvider struct {
	ID           uuid.UUID
	DisplayName  string
	Name         *string
	Image        *string
	ClientID     *string
	ClientSecret *string
}

// AuthProviderRepository defines the interface for auth provider data access
type AuthProviderRepository interface {
	// FindByName finds an auth provider by name
	FindByName(ctx context.Context, name string) (*AuthProvider, error)

	// FindByID finds an auth provider by ID
	FindByID(ctx context.Context, id uuid.UUID) (*AuthProvider, error)

	// Create creates a new auth provider
	Create(ctx context.Context, provider *AuthProvider) error
}
