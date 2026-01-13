# Local Deployment for Catetin (server-side/deployment/local) ✅

This directory contains Docker artifacts to run the Catetin server locally with a Postgres database.

What it includes:
- `Dockerfile` — multi-stage build for `./cmd/api`, copies migration files into the image so migrations are executed at app startup.
- `docker-compose.yml` — runs `db` (Postgres 15) and `app` services; the app reads from `.env.local`.
- `.env.local.example` — example env file with required values.
- `initdb/001_create_uuid_extension.sql` — creates `uuid-ossp` extension on DB initialization (only runs on first container start).

Quickstart 🔧
1. Copy the example env and edit required values (especially `DB_PASSWORD` and `JWT_SECRET_KEY`):

   cp .env.local.example .env.local
   # Edit .env.local to set secure values

> **Important:** `DB_PASSWORD` and `JWT_SECRET_KEY` are required. For Docker Compose DB initialization, set `POSTGRES_PASSWORD` in `.env.local` (or set `DB_PASSWORD` and also add `POSTGRES_PASSWORD` with the same value). If `POSTGRES_PASSWORD` is missing, the Postgres container will refuse to initialize.

2. Start services:

   cd server-side/deployment/local
   docker compose up --build

Notes & recommendations 💡
- Migrations run automatically when the app starts (this is the chosen local strategy). If the DB is not yet ready the app may exit and be restarted by Docker until the DB becomes available.
- The initdb script will create the `uuid-ossp` extension during DB initialization, which avoids permission issues where the app account lacks permission to create extensions.
- If you prefer more control over migrations, we can add a separate `migrate` one-shot service instead.

Testing & troubleshooting ⚠️
- If the app continuously restarts, check Postgres logs and confirm `.env.local` values are correct.
- To run migrations manually (inside the built image), you can run the binary with `migrate` subcommands if needed (or use `go run cmd/migrate` locally).

