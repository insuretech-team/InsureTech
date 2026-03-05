# InsureTech Backend - Database Manager

## Running DBManager

### From `backend/inscore` directory:

```powershell
# Method 1: Direct Go command
go run ./cmd/dbmanager <command> [flags]

# Method 2: Using wrapper script
./dbmanager.ps1 <command> [flags]
```

## Local Infra For AuthN/AuthZ

From repo root:

```powershell
docker compose up -d
pwsh .\scripts\generate-authn-jwt-keys.ps1
```

`docker compose up -d` starts Redis + Apache Kafka only. PostgreSQL is expected from your cloud PG config in root `.env`.

Optional: run authn inside Docker (ports from `backend/inscore/configs/services.yaml`):

```powershell
docker compose --profile authn up -d authn
```

Use env templates:

- `backend/inscore/microservices/authn/.env.example`
- `backend/inscore/microservices/authz/.env.example`

To use FLVE as external KYC provider in AuthN:

- Set `KYC_SERVICE_ENABLED=true`
- Set `KYC_SERVICE_ADDRESS=https://<your-flve-host>`
- Optionally set `KYC_FLVE_TOKEN=<bearer-token>`

### Available Commands

```powershell
# Run migrations
go run ./cmd/dbmanager migrate --target=primary
go run ./cmd/dbmanager migrate --target=backup

# Export CSV
go run ./cmd/dbmanager csv-backup --table=users --source=primary

# Run SQL query
go run ./cmd/dbmanager sql --sql="SELECT * FROM users LIMIT 10" --target=primary

# Interactive TUI mode
go run ./cmd/dbmanager
```

### Examples

```powershell
# Navigate to backend/inscore
cd E:\Projects\InsureTech\backend\inscore

# Run full migration
go run ./cmd/dbmanager migrate --target=primary

# Export users table
go run ./cmd/dbmanager csv-backup --table=users --source=primary

# Check migration status
go run ./cmd/dbmanager sql --sql="SELECT * FROM schema_migrations ORDER BY applied_at DESC LIMIT 10" --target=primary
```

### Building Standalone Binary

```powershell
cd cmd/dbmanager
go build -o dbmanager.exe .

# Then run from anywhere
./dbmanager.exe migrate --target=primary
```

## Directory Structure

```
backend/inscore/
├── cmd/
│   └── dbmanager/          # Database manager CLI
│       ├── main.go         # Entry point
│       └── internal/       # Internal packages
├── db/
│   ├── migrations/         # SQL migration files (organized by schema)
│   ├── seeds/              # Seed data files
│   ├── ops/
│   │   └── migrate.go      # Unified migration manager
│   └── config.go           # Database configuration
└── dbmanager.ps1           # Convenience wrapper
```

## Migration System

### Three-Phase Migration Flow

1. **Proto-Driven (Phase 1)**: Creates tables from proto definitions
2. **SQL Migrations (Phase 2)**: Custom SQL enhancements
3. **Seeders (Phase 3)**: Populate test/dev data

### Adding Migrations

Create files in: `db/migrations/{schema}/{timestamp}_{description}.up.sql`

Example: `db/migrations/public/20250129_001_add_user_indexes.up.sql`

### Adding Seeders

Create files in: `db/seeds/{schema}/seed_{description}.sql`

Example: `db/seeds/public/seed_test_users.sql`
