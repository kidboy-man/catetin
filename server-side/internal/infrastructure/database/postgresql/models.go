package postgresql

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// JSONB type for PostgreSQL JSONB columns
type JSONB []string

// Scan implements the sql.Scanner interface
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = []string{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to unmarshal JSONB value")
	}

	return json.Unmarshal(bytes, j)
}

// Value implements the driver.Valuer interface
func (j JSONB) Value() (driver.Value, error) {
	if len(j) == 0 {
		return json.Marshal([]string{})
	}
	return json.Marshal(j)
}

// UserModel represents the users table
type UserModel struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	FullName    string         `gorm:"type:varchar;not null"`
	PhoneNumber string         `gorm:"type:varchar;uniqueIndex;not null"`
	Image       *string        `gorm:"type:varchar"`
	Version     int            `gorm:"type:integer;not null;default:0"`
	CreatedAt   time.Time      `gorm:"type:timestamptz"`
	UpdatedAt   time.Time      `gorm:"type:timestamptz"`
	DeletedAt   gorm.DeletedAt `gorm:"type:timestamptz;index"`
}

// TableName specifies the table name for UserModel
func (UserModel) TableName() string {
	return "users"
}

// AuthProviderModel represents the auth_providers table
type AuthProviderModel struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	DisplayName  string         `gorm:"type:varchar;not null"`
	Name         *string        `gorm:"type:varchar;uniqueIndex"`
	Image        *string        `gorm:"type:varchar"`
	ClientID     *string        `gorm:"type:varchar"`
	ClientSecret *string        `gorm:"type:varchar"`
	Version      int            `gorm:"type:integer;not null;default:0"`
	CreatedAt    time.Time      `gorm:"type:timestamptz"`
	UpdatedAt    time.Time      `gorm:"type:timestamptz"`
	DeletedAt    gorm.DeletedAt `gorm:"type:timestamptz;index"`
}

// TableName specifies the table name for AuthProviderModel
func (AuthProviderModel) TableName() string {
	return "auth_providers"
}

// UserAuthModel represents the user_auths table
type UserAuthModel struct {
	ID                 uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID             uuid.UUID      `gorm:"type:uuid;not null;index:idx_user_auth_provider"`
	AuthProviderID     uuid.UUID      `gorm:"type:uuid;not null;index:idx_user_auth_provider"`
	CredentialID       string         `gorm:"type:varchar;not null"`
	CredentialSecret   string         `gorm:"type:varchar;not null"`
	CredentialRefresh  *string        `gorm:"type:varchar"`
	Version            int            `gorm:"type:integer;not null;default:0"`
	CreatedAt          time.Time      `gorm:"type:timestamptz"`
	UpdatedAt          time.Time      `gorm:"type:timestamptz"`
	DeletedAt          gorm.DeletedAt `gorm:"type:timestamptz;index"`

	// Foreign key relationships
	User         UserModel         `gorm:"foreignKey:UserID;references:ID"`
	AuthProvider AuthProviderModel `gorm:"foreignKey:AuthProviderID;references:ID"`
}

// TableName specifies the table name for UserAuthModel
func (UserAuthModel) TableName() string {
	return "user_auths"
}

// MoneyFlowModel represents the money_flows table
type MoneyFlowModel struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null;index"`
	Category    *string        `gorm:"type:varchar"`
	Amount      float64        `gorm:"type:decimal;not null"`
	Currency    string         `gorm:"type:varchar;not null;default:'IDR'"`
	Description *string        `gorm:"type:text"`
	Tags        JSONB          `gorm:"type:jsonb"`
	Version     int            `gorm:"type:integer;not null;default:0"`
	CreatedAt   time.Time      `gorm:"type:timestamptz"`
	UpdatedAt   time.Time      `gorm:"type:timestamptz"`
	DeletedAt   gorm.DeletedAt `gorm:"type:timestamptz;index"`

	// Foreign key relationship
	User UserModel `gorm:"foreignKey:UserID;references:ID"`
}

// TableName specifies the table name for MoneyFlowModel
func (MoneyFlowModel) TableName() string {
	return "money_flows"
}
