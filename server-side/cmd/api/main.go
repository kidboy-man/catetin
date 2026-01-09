package main

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/ingunawandra/catetin/internal/config"
	httpController "github.com/ingunawandra/catetin/internal/controller/http"
	v1 "github.com/ingunawandra/catetin/internal/controller/http/v1"
	"github.com/ingunawandra/catetin/internal/infrastructure/database/postgresql"
	"github.com/ingunawandra/catetin/internal/infrastructure/security"
	"github.com/ingunawandra/catetin/internal/service"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	log.Printf("Starting Catetin API Server on port %s...", cfg.Server.Port)

	// Initialize database connection
	db, err := postgresql.NewConnection(cfg.GetDatabaseDSN(), cfg.Server.Env)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run database migrations using golang-migrate
	databaseURL, err := postgresql.ConvertDSNToURL(cfg.GetDatabaseDSN())
	if err != nil {
		log.Fatalf("Failed to convert DSN to URL: %v", err)
	}

	// Get absolute path to migrations directory
	migrationsPath, err := filepath.Abs("internal/infrastructure/database/postgresql/migrations")
	if err != nil {
		log.Fatalf("Failed to get migrations path: %v", err)
	}

	// Run migrations
	if err := postgresql.RunMigrations(databaseURL, migrationsPath); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}

	// Check migration version
	version, dirty, err := postgresql.MigrationVersion(databaseURL, migrationsPath)
	if err != nil {
		log.Printf("Warning: Failed to get migration version: %v", err)
	} else {
		log.Printf("Current database migration version: %d (dirty: %v)", version, dirty)
	}

	// Initialize repositories
	userRepo := postgresql.NewUserRepository(db)
	moneyFlowRepo := postgresql.NewMoneyFlowRepository(db)
	authProviderRepo := postgresql.NewAuthProviderRepository(db)
	userAuthRepo := postgresql.NewUserAuthRepository(db)

	_ = moneyFlowRepo

	// Initialize transaction manager
	txManager := postgresql.NewTransactionManager(db)

	// Initialize security utilities
	passwordHasher := security.NewPasswordHasher()
	jwtManager := security.NewJWTManager(
		cfg.JWT.SecretKey,
		time.Duration(cfg.JWT.AccessTokenDuration)*time.Minute,
		time.Duration(cfg.JWT.RefreshTokenDuration)*24*time.Hour,
	)

	// Initialize services
	authService := service.NewAuthService(
		userRepo,
		userAuthRepo,
		authProviderRepo,
		passwordHasher,
		jwtManager,
		txManager,
	)

	// Ensure email-password auth provider exists
	ctx := context.Background()
	if err := authService.EnsureEmailPasswordProvider(ctx); err != nil {
		log.Fatalf("Failed to ensure email-password auth provider: %v", err)
	}
	log.Println("Email-password authentication provider initialized")

	// Initialize HTTP handlers
	authHandler := v1.NewAuthHandler(authService)

	// Setup router
	router := httpController.SetupRouter(&httpController.RouterConfig{
		AuthHandler: authHandler,
	})

	// Start HTTP server
	serverAddr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Starting HTTP server on %s...", serverAddr)
	log.Println("Phase 1: Project Setup & Database Foundation - COMPLETED")
	log.Println("Authentication endpoints available:")
	log.Println("  POST /api/v1/authentications/register")
	log.Println("  POST /api/v1/authentications/login")
	log.Println("  GET  /health")

	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}
