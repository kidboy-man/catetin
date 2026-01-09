-- Drop indexes first (optional but clean)
DROP INDEX IF EXISTS idx_money_flows_category;
DROP INDEX IF EXISTS idx_money_flows_deleted_at;
DROP INDEX IF EXISTS idx_money_flows_created_at;
DROP INDEX IF EXISTS idx_money_flows_user_id;

DROP INDEX IF EXISTS idx_user_auths_deleted_at;
DROP INDEX IF EXISTS idx_user_auths_user_provider;

DROP INDEX IF EXISTS idx_auth_providers_deleted_at;
DROP INDEX IF EXISTS idx_auth_providers_name_unique;

DROP INDEX IF EXISTS idx_users_deleted_at;
DROP INDEX IF EXISTS idx_users_phone_number;
DROP INDEX IF EXISTS idx_users_phone_number_unique;

-- Drop tables in reverse order (respecting foreign key constraints)
DROP TABLE IF EXISTS "money_flows" CASCADE;
DROP TABLE IF EXISTS "user_auths" CASCADE;
DROP TABLE IF EXISTS "auth_providers" CASCADE;
DROP TABLE IF EXISTS "users" CASCADE;

-- Drop UUID extension (optional, keep it if other tables use it)
-- DROP EXTENSION IF EXISTS "uuid-ossp";
