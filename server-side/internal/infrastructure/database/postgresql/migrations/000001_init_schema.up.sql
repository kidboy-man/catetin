-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create users table
CREATE TABLE IF NOT EXISTS "users" (
  "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  "full_name" varchar NOT NULL,
  "phone_number" varchar NOT NULL,
  "image" varchar,
  "version" integer NOT NULL DEFAULT 0,
  "created_at" timestamptz NOT NULL DEFAULT NOW(),
  "updated_at" timestamptz NOT NULL DEFAULT NOW(),
  "deleted_at" timestamptz
);

-- Create unique index on phone_number only for non-deleted records
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_phone_number_unique ON "users" ("phone_number") WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_users_phone_number ON "users" ("phone_number");
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON "users" ("deleted_at");

-- Create auth_providers table
CREATE TABLE IF NOT EXISTS "auth_providers" (
  "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  "display_name" varchar NOT NULL,
  "name" varchar,
  "image" varchar,
  "client_id" varchar,
  "client_secret" varchar,
  "version" integer NOT NULL DEFAULT 0,
  "created_at" timestamptz NOT NULL DEFAULT NOW(),
  "updated_at" timestamptz NOT NULL DEFAULT NOW(),
  "deleted_at" timestamptz
);

-- Create unique index on name only for non-deleted records
CREATE UNIQUE INDEX IF NOT EXISTS idx_auth_providers_name_unique ON "auth_providers" ("name") WHERE deleted_at IS NULL AND name IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_auth_providers_deleted_at ON "auth_providers" ("deleted_at");

-- Create user_auths table
CREATE TABLE IF NOT EXISTS "user_auths" (
  "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  "user_id" uuid NOT NULL,
  "auth_provider_id" uuid NOT NULL,
  "credential_id" varchar NOT NULL,
  "credential_secret" varchar NOT NULL,
  "credential_refresh" varchar,
  "version" integer NOT NULL DEFAULT 0,
  "created_at" timestamptz NOT NULL DEFAULT NOW(),
  "updated_at" timestamptz NOT NULL DEFAULT NOW(),
  "deleted_at" timestamptz,
  CONSTRAINT fk_user_auths_user FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE,
  CONSTRAINT fk_user_auths_auth_provider FOREIGN KEY ("auth_provider_id") REFERENCES "auth_providers" ("id") ON DELETE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_user_auths_user_provider ON "user_auths" ("user_id", "auth_provider_id") WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_user_auths_deleted_at ON "user_auths" ("deleted_at");

-- Create money_flows table
CREATE TABLE IF NOT EXISTS "money_flows" (
  "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  "user_id" uuid NOT NULL,
  "category" varchar,
  "amount" decimal NOT NULL,
  "currency" varchar NOT NULL DEFAULT 'IDR',
  "description" text,
  "tags" jsonb DEFAULT '[]'::jsonb,
  "version" integer NOT NULL DEFAULT 0,
  "created_at" timestamptz NOT NULL DEFAULT NOW(),
  "updated_at" timestamptz NOT NULL DEFAULT NOW(),
  "deleted_at" timestamptz,
  CONSTRAINT fk_money_flows_user FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_money_flows_user_id ON "money_flows" ("user_id");
CREATE INDEX IF NOT EXISTS idx_money_flows_created_at ON "money_flows" ("created_at");
CREATE INDEX IF NOT EXISTS idx_money_flows_deleted_at ON "money_flows" ("deleted_at");
CREATE INDEX IF NOT EXISTS idx_money_flows_category ON "money_flows" ("category");

-- Add comments for documentation
COMMENT ON TABLE "users" IS 'User profiles and authentication information';
COMMENT ON TABLE "auth_providers" IS 'OAuth and authentication provider configurations';
COMMENT ON TABLE "user_auths" IS 'User authentication credentials linked to auth providers';
COMMENT ON TABLE "money_flows" IS 'User expense and money flow tracking records';

COMMENT ON COLUMN "users"."phone_number" IS 'Unique phone number for WhatsApp integration';
COMMENT ON COLUMN "money_flows"."tags" IS 'JSONB array of tags for categorization and filtering';
COMMENT ON COLUMN "users"."version" IS 'Version field for optimistic locking';
COMMENT ON COLUMN "money_flows"."version" IS 'Version field for optimistic locking';
