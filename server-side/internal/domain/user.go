package domain

import (
	"time"

	"github.com/google/uuid"
)

// User represents the core user entity
type User struct {
	ID          uuid.UUID
	FullName    string
	PhoneNumber string
	Image       *string
	Version     int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

// NewUser creates a new User entity
func NewUser(fullName, phoneNumber string) *User {
	now := time.Now()
	return &User{
		ID:          uuid.New(),
		FullName:    fullName,
		PhoneNumber: phoneNumber,
		Version:     0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// IsDeleted checks if the user is soft deleted
func (u *User) IsDeleted() bool {
	return u.DeletedAt != nil
}

// IncrementVersion increments the version for optimistic locking
func (u *User) IncrementVersion() {
	u.Version++
	u.UpdatedAt = time.Now()
}

// SoftDelete marks the user as deleted
func (u *User) SoftDelete() {
	now := time.Now()
	u.DeletedAt = &now
	u.UpdatedAt = now
}
