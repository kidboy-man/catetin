-- This script will be executed by the official Postgres image on first initialization.
-- It creates the uuid-ossp extension required by the initial migration.

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
