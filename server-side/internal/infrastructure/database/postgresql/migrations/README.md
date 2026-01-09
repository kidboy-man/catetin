# Database Migrations

This directory contains database migration files managed by [golang-migrate](https://github.com/golang-migrate/migrate).

## Migration Files

Migration files follow the naming convention:
```
{version}_{name}.up.sql    # Applied when migrating up
{version}_{name}.down.sql  # Applied when rolling back
```

Example:
- `000001_init_schema.up.sql` - Creates initial database schema
- `000001_init_schema.down.sql` - Drops initial database schema

## Creating New Migrations

To create a new migration, create two files with the next version number:

```bash
# Example: Adding a new column
touch 000002_add_user_timezone.up.sql
touch 000002_add_user_timezone.down.sql
```

**000002_add_user_timezone.up.sql:**
```sql
ALTER TABLE users ADD COLUMN timezone VARCHAR(50) DEFAULT 'Asia/Jakarta';
```

**000002_add_user_timezone.down.sql:**
```sql
ALTER TABLE users DROP COLUMN timezone;
```

## Running Migrations

### Automatic (on application start)
Migrations run automatically when the application starts via `main.go`.

### Manual using CLI tool
You can also run migrations manually using the CLI tool:

```bash
# Run all pending migrations
go run cmd/migrate/main.go up

# Rollback last migration
go run cmd/migrate/main.go down

# Rollback N migrations
go run cmd/migrate/main.go down -steps 2

# Check current version
go run cmd/migrate/main.go version

# Force version (use with caution!)
go run cmd/migrate/main.go force -version 1
```

## Best Practices

1. **Always write both up and down migrations**
   - Up: Apply the change
   - Down: Revert the change

2. **Keep migrations small and focused**
   - One logical change per migration
   - Easier to review and rollback

3. **Test migrations in development first**
   - Apply migration: `go run cmd/migrate/main.go up`
   - Verify changes in database
   - Test rollback: `go run cmd/migrate/main.go down`
   - Re-apply: `go run cmd/migrate/main.go up`

4. **Never modify existing migration files**
   - Once applied to production, migrations are immutable
   - Create a new migration to fix issues

5. **Use transactions when possible**
   - Wrap DDL statements in transactions for safety
   - PostgreSQL supports transactional DDL

6. **Handle data migrations carefully**
   - Consider data volume
   - May need to run migrations in batches
   - Test with production-like data

## Troubleshooting

### Dirty State
If migrations fail midway, the database may be in a "dirty" state:

```bash
# Check current state
go run cmd/migrate/main.go version

# If dirty, fix manually then force version
go run cmd/migrate/main.go force -version <last_good_version>
```

### Migration Failed
If a migration fails:

1. Check error logs for the cause
2. Fix the SQL in the migration file
3. Rollback to previous version: `go run cmd/migrate/main.go down`
4. Re-apply migrations: `go run cmd/migrate/main.go up`

## PostgreSQL Connection

Migrations use the same database configuration from `.env`:
- `DB_HOST`
- `DB_PORT`
- `DB_USER`
- `DB_PASSWORD`
- `DB_NAME`
- `DB_SSLMODE`
