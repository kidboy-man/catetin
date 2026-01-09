# Database Migrations with golang-migrate

This project uses [golang-migrate](https://github.com/golang-migrate/migrate) for database schema management.

## Why golang-migrate?

- **Version Control**: Track database schema changes alongside code
- **Reproducible**: Apply same migrations across dev, staging, and production
- **Rollback Support**: Safely revert changes if needed
- **Team Collaboration**: Multiple developers can work on schema changes
- **CI/CD Ready**: Automated migration runs in deployment pipelines

## Quick Start

### 1. Set up your database

Create a PostgreSQL database and configure `.env`:

```bash
cp .env.example .env
# Edit .env with your database credentials
```

### 2. Run migrations

**Option A: Automatic (Recommended)**
Migrations run automatically when starting the application:

```bash
go run cmd/api/main.go
```

**Option B: Manual using CLI tool**

```bash
# Apply all pending migrations
go run cmd/migrate/main.go up

# Check current version
go run cmd/migrate/main.go version
```

## Migration Files Location

All migration files are stored in:
```
internal/infrastructure/database/postgresql/migrations/
```

## Current Migrations

### 000001_init_schema
Creates the initial database schema with:
- `users` table - User profiles with phone numbers for WhatsApp
- `auth_providers` table - OAuth provider configurations
- `user_auths` table - User authentication credentials
- `money_flows` table - Expense tracking with JSONB tags

Features:
- UUID primary keys with auto-generation
- Soft delete support (deleted_at timestamps)
- Optimistic locking (version integers)
- Proper foreign key constraints with CASCADE delete
- Indexes on frequently queried columns
- JSONB support for flexible tags

## Creating New Migrations

### Step 1: Create migration files

Create both UP and DOWN migration files:

```bash
# Replace 000002 with next sequence number
touch internal/infrastructure/database/postgresql/migrations/000002_add_feature.up.sql
touch internal/infrastructure/database/postgresql/migrations/000002_add_feature.down.sql
```

### Step 2: Write SQL

**000002_add_feature.up.sql** (what to apply):
```sql
ALTER TABLE users ADD COLUMN timezone VARCHAR(50) DEFAULT 'Asia/Jakarta';
CREATE INDEX idx_users_timezone ON users(timezone);
```

**000002_add_feature.down.sql** (how to rollback):
```sql
DROP INDEX IF EXISTS idx_users_timezone;
ALTER TABLE users DROP COLUMN timezone;
```

### Step 3: Test migrations

```bash
# Apply new migration
go run cmd/migrate/main.go up

# Verify in database
# psql -U postgres -d catetin -c "\d users"

# Test rollback
go run cmd/migrate/main.go down

# Re-apply
go run cmd/migrate/main.go up
```

## CLI Commands

### Migrate Up
Apply all pending migrations:
```bash
go run cmd/migrate/main.go up
```

### Migrate Down
Rollback the last migration:
```bash
go run cmd/migrate/main.go down
```

Rollback multiple migrations:
```bash
go run cmd/migrate/main.go down -steps 3
```

### Check Version
Show current migration version:
```bash
go run cmd/migrate/main.go version
```

### Force Version (Use with caution!)
If database is in dirty state:
```bash
go run cmd/migrate/main.go force -version 1
```

## Migration Best Practices

### ✅ DO:

1. **Always create both UP and DOWN migrations**
   - UP applies the change
   - DOWN reverts it completely

2. **Keep migrations atomic and focused**
   - One logical change per migration
   - Easier to understand and rollback

3. **Test thoroughly before production**
   ```bash
   # Test cycle
   go run cmd/migrate/main.go up      # Apply
   # Verify changes work
   go run cmd/migrate/main.go down    # Rollback
   # Verify rollback works
   go run cmd/migrate/main.go up      # Re-apply
   ```

4. **Use meaningful names**
   - Good: `000002_add_user_timezone.sql`
   - Bad: `000002_update.sql`

5. **Add comments in complex migrations**
   ```sql
   -- Add timezone support for international users
   ALTER TABLE users ADD COLUMN timezone VARCHAR(50);
   ```

### ❌ DON'T:

1. **Never modify existing migration files**
   - Once applied to production, treat as immutable
   - Create a new migration to fix issues

2. **Don't skip version numbers**
   - Maintain sequential order: 000001, 000002, 000003...

3. **Don't forget indexes**
   - Add indexes for foreign keys and frequently queried columns

4. **Don't ignore DOWN migrations**
   - Always write proper rollback logic

## Troubleshooting

### Dirty State Error

If migrations fail midway, the database enters "dirty" state:

```bash
# Check status
go run cmd/migrate/main.go version
# Output: Current version: 1 (DIRTY)

# Fix the issue, then force to last good version
go run cmd/migrate/main.go force -version 1

# Try again
go run cmd/migrate/main.go up
```

### Migration Fails

1. Check error logs for SQL syntax errors
2. Fix the migration file
3. Rollback: `go run cmd/migrate/main.go down`
4. Re-apply: `go run cmd/migrate/main.go up`

### Cannot Find Migrations

Ensure you're running commands from the project root:
```bash
cd /path/to/catetin/server-side
go run cmd/migrate/main.go up
```

## Integration with Application

The application automatically runs migrations on startup in [cmd/api/main.go](cmd/api/main.go:26-41):

```go
// Run database migrations using golang-migrate
databaseURL, err := postgresql.ConvertDSNToURL(cfg.GetDatabaseDSN())
if err != nil {
    log.Fatalf("Failed to convert DSN to URL: %v", err)
}

migrationsPath, err := filepath.Abs("internal/infrastructure/database/postgresql/migrations")
if err != nil {
    log.Fatalf("Failed to get migrations path: %v", err)
}

if err := postgresql.RunMigrations(databaseURL, migrationsPath); err != nil {
    log.Fatalf("Failed to run database migrations: %v", err)
}
```

## Production Deployment

### Option 1: Automatic on startup (Current)
Migrations run when the application starts. Simple but requires careful testing.

### Option 2: Separate migration step (Recommended for production)
Run migrations as a separate step in your deployment pipeline:

```bash
# In your CI/CD pipeline
go run cmd/migrate/main.go up
go run cmd/api/main.go
```

### Option 3: Use migrate CLI directly
Install the migrate CLI tool:
```bash
# macOS
brew install golang-migrate

# Linux
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.1/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/

# Run migrations
migrate -path internal/infrastructure/database/postgresql/migrations -database "postgres://user:pass@localhost:5432/catetin?sslmode=disable" up
```

## Further Reading

- [golang-migrate Documentation](https://github.com/golang-migrate/migrate)
- [PostgreSQL ALTER TABLE](https://www.postgresql.org/docs/current/sql-altertable.html)
- [Migration Best Practices](https://github.com/golang-migrate/migrate/blob/master/MIGRATIONS.md)
