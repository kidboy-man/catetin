package domain

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

// MoneyFlow represents the core expense/money flow entity
type MoneyFlow struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Category    *string
	Amount      float64
	Currency    string
	Description *string
	Tags        []string
	Version     int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

// NewMoneyFlow creates a new MoneyFlow entity
func NewMoneyFlow(userID uuid.UUID, amount float64, currency string) (*MoneyFlow, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}

	if currency == "" {
		currency = "IDR" // Default to Indonesian Rupiah
	}

	now := time.Now()
	return &MoneyFlow{
		ID:        uuid.New(),
		UserID:    userID,
		Amount:    amount,
		Currency:  currency,
		Version:   0,
		CreatedAt: now,
		UpdatedAt: now,
		Tags:      []string{},
	}, nil
}

// SetCategory sets the category for the money flow
func (mf *MoneyFlow) SetCategory(category string) {
	mf.Category = &category
	mf.UpdatedAt = time.Now()
}

// SetDescription sets the description for the money flow
func (mf *MoneyFlow) SetDescription(description string) {
	mf.Description = &description
	mf.UpdatedAt = time.Now()
}

// AddTag adds a tag to the money flow
func (mf *MoneyFlow) AddTag(tag string) {
	if mf.Tags == nil {
		mf.Tags = []string{}
	}
	mf.Tags = append(mf.Tags, tag)
	mf.UpdatedAt = time.Now()
}

// SetTags sets all tags for the money flow
func (mf *MoneyFlow) SetTags(tags []string) {
	mf.Tags = tags
	mf.UpdatedAt = time.Now()
}

// IsDeleted checks if the money flow is soft deleted
func (mf *MoneyFlow) IsDeleted() bool {
	return mf.DeletedAt != nil
}

// IncrementVersion increments the version for optimistic locking
func (mf *MoneyFlow) IncrementVersion() {
	mf.Version++
	mf.UpdatedAt = time.Now()
}

// SoftDelete marks the money flow as deleted
func (mf *MoneyFlow) SoftDelete() {
	now := time.Now()
	mf.DeletedAt = &now
	mf.UpdatedAt = now
}

// TagsToJSON converts tags slice to JSON for database storage
func (mf *MoneyFlow) TagsToJSON() ([]byte, error) {
	if len(mf.Tags) == 0 {
		return json.Marshal([]string{})
	}
	return json.Marshal(mf.Tags)
}

// TagsFromJSON populates tags from JSON data
func (mf *MoneyFlow) TagsFromJSON(data []byte) error {
	if len(data) == 0 {
		mf.Tags = []string{}
		return nil
	}
	return json.Unmarshal(data, &mf.Tags)
}
