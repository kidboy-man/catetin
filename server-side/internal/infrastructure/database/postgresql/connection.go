package postgresql

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewConnection creates a new PostgreSQL database connection
func NewConnection(dsn string, env string) (*gorm.DB, error) {
	// Configure GORM logger based on environment
	logLevel := logger.Info
	if env == "production" {
		logLevel = logger.Warn
	}

	config := &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}

	// Open connection
	db, err := gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying SQL database
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL database")

	return db, nil
}

// AutoMigrate runs GORM auto-migration for all models
// NOTE: This is deprecated in favor of golang-migrate. Use only for development/testing.
func AutoMigrate(db *gorm.DB) error {
	log.Println("Running GORM auto-migrations (deprecated - use golang-migrate instead)...")

	err := db.AutoMigrate(
		&UserModel{},
		&MoneyFlowModel{},
		&AuthProviderModel{},
		&UserAuthModel{},
	)

	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("GORM auto-migrations completed successfully")
	return nil
}

// ConvertDSNToURL converts GORM DSN format to standard PostgreSQL URL format for golang-migrate
// Example: "host=localhost port=5432 user=postgres password=pass dbname=catetin sslmode=disable"
// becomes: "postgres://postgres:pass@localhost:5432/catetin?sslmode=disable"
func ConvertDSNToURL(dsn string) (string, error) {
	// Parse DSN manually
	params := make(map[string]string)
	var key, value string
	var inValue bool

	for i := 0; i < len(dsn); i++ {
		char := dsn[i]
		if char == '=' {
			inValue = true
			continue
		}
		if char == ' ' {
			if inValue && key != "" {
				params[key] = value
				key = ""
				value = ""
				inValue = false
			}
			continue
		}
		if inValue {
			value += string(char)
		} else {
			key += string(char)
		}
	}
	if key != "" && value != "" {
		params[key] = value
	}

	// Extract required parameters
	host := params["host"]
	port := params["port"]
	user := params["user"]
	password := params["password"]
	dbname := params["dbname"]
	sslmode := params["sslmode"]

	if host == "" || port == "" || user == "" || dbname == "" {
		return "", fmt.Errorf("invalid DSN format: missing required parameters")
	}

	// Construct PostgreSQL URL
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, dbname)
	if sslmode != "" {
		url += fmt.Sprintf("?sslmode=%s", sslmode)
	}

	return url, nil
}
