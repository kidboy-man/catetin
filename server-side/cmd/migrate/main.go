package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ingunawandra/catetin/internal/config"
	"github.com/ingunawandra/catetin/internal/infrastructure/database/postgresql"
)

func main() {
	// Define subcommands
	upCmd := flag.NewFlagSet("up", flag.ExitOnError)
	downCmd := flag.NewFlagSet("down", flag.ExitOnError)
	versionCmd := flag.NewFlagSet("version", flag.ExitOnError)
	forceCmd := flag.NewFlagSet("force", flag.ExitOnError)

	// Flags for down command
	downSteps := downCmd.Int("steps", 1, "Number of migrations to rollback")

	// Flags for force command
	forceVersion := forceCmd.Int("version", -1, "Version to force")

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Convert DSN to URL
	databaseURL, err := postgresql.ConvertDSNToURL(cfg.GetDatabaseDSN())
	if err != nil {
		log.Fatalf("Failed to convert DSN to URL: %v", err)
	}

	// Get migrations path
	migrationsPath, err := filepath.Abs("internal/infrastructure/database/postgresql/migrations")
	if err != nil {
		log.Fatalf("Failed to get migrations path: %v", err)
	}

	// Parse subcommand
	switch os.Args[1] {
	case "up":
		upCmd.Parse(os.Args[2:])
		if err := postgresql.RunMigrations(databaseURL, migrationsPath); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		fmt.Println("✅ All migrations applied successfully")

	case "down":
		downCmd.Parse(os.Args[2:])
		if err := postgresql.RollbackMigration(databaseURL, migrationsPath, *downSteps); err != nil {
			log.Fatalf("Rollback failed: %v", err)
		}
		fmt.Printf("✅ Successfully rolled back %d migration(s)\n", *downSteps)

	case "version":
		versionCmd.Parse(os.Args[2:])
		version, dirty, err := postgresql.MigrationVersion(databaseURL, migrationsPath)
		if err != nil {
			log.Fatalf("Failed to get version: %v", err)
		}
		if dirty {
			fmt.Printf("⚠️  Current version: %d (DIRTY - needs manual intervention)\n", version)
		} else {
			fmt.Printf("✅ Current version: %d\n", version)
		}

	case "force":
		forceCmd.Parse(os.Args[2:])
		if *forceVersion < 0 {
			log.Fatal("Please specify a version using -version flag")
		}
		if err := postgresql.ForceMigrationVersion(databaseURL, migrationsPath, *forceVersion); err != nil {
			log.Fatalf("Force version failed: %v", err)
		}
		fmt.Printf("✅ Forced version to %d\n", *forceVersion)

	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Database Migration Tool")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  go run cmd/migrate/main.go <command> [options]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  up                    Apply all pending migrations")
	fmt.Println("  down [-steps N]       Rollback N migrations (default: 1)")
	fmt.Println("  version               Show current migration version")
	fmt.Println("  force -version N      Force migration version (use with caution!)")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  go run cmd/migrate/main.go up")
	fmt.Println("  go run cmd/migrate/main.go down")
	fmt.Println("  go run cmd/migrate/main.go down -steps 2")
	fmt.Println("  go run cmd/migrate/main.go version")
	fmt.Println("  go run cmd/migrate/main.go force -version 1")
}
