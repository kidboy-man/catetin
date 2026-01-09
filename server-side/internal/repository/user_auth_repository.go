package repository

import (
	"context"

	"github.com/google/uuid"
)

// UserAuth represents a user authentication record
type UserAuth struct {
	ID                uuid.UUID
	UserID            uuid.UUID
	AuthProviderID    uuid.UUID
	CredentialID      string // email for email-password auth
	CredentialSecret  string // hashed password
	CredentialRefresh *string
}

// UserAuthRepository defines the interface for user auth data access
type UserAuthRepository interface {
	// Create creates a new user auth record
	Create(ctx context.Context, userAuth *UserAuth) error

	// FindByCredentialID finds a user auth by credential ID (email) and provider
	FindByCredentialID(ctx context.Context, credentialID string, authProviderID uuid.UUID) (*UserAuth, error)

	// FindByUserIDAndProvider finds a user auth by user ID and provider
	FindByUserIDAndProvider(ctx context.Context, userID, authProviderID uuid.UUID) (*UserAuth, error)

	// Update updates a user auth record
	Update(ctx context.Context, userAuth *UserAuth) error

	// Delete soft deletes a user auth record
	Delete(ctx context.Context, id uuid.UUID) error
}
