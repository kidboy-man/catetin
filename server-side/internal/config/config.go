package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Database  DatabaseConfig
	OpenAI    OpenAIConfig
	WhatsApp  WhatsAppConfig
	Server    ServerConfig
	Webhook   WebhookConfig
	JWT       JWTConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type OpenAIConfig struct {
	APIKey string
	Model  string
}

type WhatsAppConfig struct {
	PhoneNumberID       string
	BusinessAccountID   string
	AccessToken         string
	APIVersion          string
}

type ServerConfig struct {
	Port string
	Env  string
}

type WebhookConfig struct {
	VerifyToken string
}

type JWTConfig struct {
	SecretKey            string
	AccessTokenDuration  int // in minutes
	RefreshTokenDuration int // in days
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists (for local development)
	_ = godotenv.Load()

	config := &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "catetin"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		OpenAI: OpenAIConfig{
			APIKey: getEnv("OPENAI_API_KEY", ""),
			Model:  getEnv("OPENAI_MODEL", "gpt-4o-mini"),
		},
		WhatsApp: WhatsAppConfig{
			PhoneNumberID:     getEnv("WHATSAPP_PHONE_NUMBER_ID", ""),
			BusinessAccountID: getEnv("WHATSAPP_BUSINESS_ACCOUNT_ID", ""),
			AccessToken:       getEnv("WHATSAPP_ACCESS_TOKEN", ""),
			APIVersion:        getEnv("WHATSAPP_API_VERSION", "v21.0"),
		},
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Env:  getEnv("ENV", "development"),
		},
		Webhook: WebhookConfig{
			VerifyToken: getEnv("WEBHOOK_VERIFY_TOKEN", ""),
		},
		JWT: JWTConfig{
			SecretKey:            getEnv("JWT_SECRET_KEY", ""),
			AccessTokenDuration:  getEnvAsInt("JWT_ACCESS_TOKEN_DURATION", 60),   // 60 minutes default
			RefreshTokenDuration: getEnvAsInt("JWT_REFRESH_TOKEN_DURATION", 30), // 30 days default
		},
	}

	// Validate required fields
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Database.Password == "" {
		return fmt.Errorf("DB_PASSWORD is required")
	}

	if c.JWT.SecretKey == "" {
		return fmt.Errorf("JWT_SECRET_KEY is required")
	}

	// Note: OpenAI, WhatsApp, and Webhook configs are optional
	// They will be validated when those features are used

	return nil
}

// GetDatabaseDSN returns the PostgreSQL connection string
func (c *Config) GetDatabaseDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	var value int
	_, err := fmt.Sscanf(valueStr, "%d", &value)
	if err != nil {
		return defaultValue
	}
	return value
}
