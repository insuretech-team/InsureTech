# InsureTech Database Migration System

## 🎯 Overview

This is a **unified, proto-driven migration system** that automatically discovers schemas from proto files and executes migrations in three phases:

1. **Proto-driven table creation** (Phase 1)
2. **Custom SQL migrations** (Phase 2)  
3. **Seed data** (Phase 3)

---

## 📁 Directory Structure

```
backend/inscore/db/
├── migrations/              # Custom SQL migrations (organized by schema)
│   ├── public/
│   │   └── 20250101_001_example.up.sql
│   ├── policy_schema/
│   ├── analytics_schema/
│   └── [other schemas]/
├── seeds/                   # Seed data files (organized by schema)
│   ├── public/
│   │   └── seed_users.sql
│   └── [other schemas]/
└── ops/
    └── migrate.go          # Unified migration manager

proto/insuretech/           # Proto definitions with table annotations
├── authn/entity/v1/
├── policy/entity/v1/
└── [30+ domain folders]/
```

---

## 🚀 Quick Start

### Run All Migrations

```bash
# From project root
cd backend/inscore/cmd/dbmanager
go run main.go migrate --target=primary
```

### What Happens:

```
==========================================
🔄 Starting Unified Migration Flow
==========================================

📦 PHASE 1: Proto-Driven Table Creation
------------------------------------------
✓ Discovered 10 schemas from proto files
✓ Created 58 tables in migration order
✓ Added foreign key constraints
✓ Added indexes and comments

🔧 PHASE 2: Custom SQL Migrations
------------------------------------------
✓ Applied migration: public/20250101_001_example.up.sql
✓ Applied migration: policy_schema/20250102_001_triggers.up.sql

🌱 PHASE 3: Seeding Data
------------------------------------------
✓ Applied seeder: public/seed_users.sql

==========================================
✅ All Migrations Completed Successfully!
==========================================
```

---

## 🛡️ Safety & Integrity Flags (Declarative Mode)

The migration system supports strict, declarative synchronization:

### `--prune`
Automatically drops columns that exist in the database but are NOT defined in Proto.

```bash
dbmanager migrate --target=primary --prune
```

**Warning:** This is a destructive operation. Use with caution.

### `--strict`
Fails the migration if any schema drift is detected (zombie columns, type mismatches).

```bash
dbmanager migrate --target=primary --strict
```

**Use Case:** CI/CD pipelines to ensure database schema matches Proto definitions exactly.

### Combined Usage
```bash
# First run strict to check for drift
dbmanager migrate --target=primary --strict

# If drift found, run prune to fix it
dbmanager migrate --target=primary --prune
```

---

## 📋 Discovered Schemas

The system automatically discovers these schemas from proto files:

- **public** (50 tables) - Core entities: users, sessions, roles, etc.
- **ai_schema** (1 table) - AI agents
- **analytics_schema** (2 tables) - Business metrics
- **claims_schema** (1 table) - Insurance claims
- **iot_schema** (1 table) - IoT devices
- **notification_schema** (1 table) - Notifications
- **partner_schema** (1 table) - Partners
- **payment_schema** (1 table) - Payments
- **policy_schema** (1 table) - Insurance policies
- **product_schema** (1 table) - Insurance products

---

## 🎨 Proto Table Definition

Example from `proto/insuretech/authn/entity/v1/user.proto`:

```protobuf
message User {
  option (insuretech.common.v1.table) = {
    table_name: "users"
    schema_name: "public"
    migration_order: 1
    is_table: true
    comment: "Registered users with authentication"
    soft_delete: true
    audit_fields: true
    enable_rls: true
  };

  string id = 1 [
    (insuretech.common.v1.column) = {
      sql_type: "UUID"
      primary_key: true
      default_value: "uuid_generate_v4()"
      comment: "Unique user identifier"
    }
  ];

  string email = 2 [
    (insuretech.common.v1.column) = {
      not_null: true
      unique: true
      comment: "User email address"
    }
  ];
}
```

---

## 📝 Creating Custom Migrations

### File Naming Convention

```
{timestamp}_{description}.up.sql
```

### Example Migration

**File:** `backend/inscore/db/migrations/public/20250101_001_add_email_verified.up.sql`

```sql
-- Add email verification column to users table
-- Idempotent: checks before adding

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_schema = 'public' 
        AND table_name = 'users' 
        AND column_name = 'email_verified'
    ) THEN
        ALTER TABLE public.users 
        ADD COLUMN email_verified BOOLEAN DEFAULT false;
        
        COMMENT ON COLUMN public.users.email_verified 
        IS 'Whether user email has been verified';
    END IF;
END $$;
```

### Best Practices

1. **Idempotent**: Always check before creating/altering
2. **Transactional**: Each migration runs in a transaction
3. **Documented**: Add comments explaining the change
4. **Schema-specific**: Place in correct schema folder

---

## 🌱 Creating Seeders

### File Naming Convention

```
seed_{description}.sql
```

### Example Seeder

**File:** `backend/inscore/db/seeds/public/seed_test_users.sql`

```sql
-- Insert test users for development
-- Idempotent: only inserts if not exists

INSERT INTO public.users (id, email, password_hash, created_at, updated_at)
SELECT 
    'a0000000-0000-0000-0000-000000000001'::uuid,
    'admin@example.com',
    '$2a$10$EXAMPLE_HASH',
    now(),
    now()
WHERE NOT EXISTS (
    SELECT 1 FROM public.users WHERE email = 'admin@example.com'
);

INSERT INTO public.users (id, email, password_hash, created_at, updated_at)
SELECT 
    'a0000000-0000-0000-0000-000000000002'::uuid,
    'user@example.com',
    '$2a$10$EXAMPLE_HASH',
    now(),
    now()
WHERE NOT EXISTS (
    SELECT 1 FROM public.users WHERE email = 'user@example.com'
);
```

---

## 🔍 Migration Tracking

All migrations are tracked in the `public.schema_migrations` table:

```sql
SELECT * FROM public.schema_migrations ORDER BY applied_at DESC;
```

**Columns:**
- `id` - UUID
- `name` - Migration/seeder file name
- `type` - 'proto', 'migration', or 'seeder'
- `schema` - Target schema
- `applied_at` - Timestamp
- `checksum` - SHA256 hash for integrity
- `execution_ms` - Execution time
- `status` - 'success' or 'failed'
- `error_msg` - Error details if failed

---

## 🛠️ Advanced Usage

### Check Migration Status

```bash
go run main.go migrate --target=primary
# View the schema_migrations table
```

### Migrate Backup Database

```bash
go run main.go migrate --target=backup
```

### Add New Schema

1. Create proto file with `schema_name` annotation
2. Run migrations - schema auto-discovered
3. Directory created automatically

---

## 🔥 Key Features

✅ **Auto-discovery**: Schemas discovered from proto files  
✅ **Three-phase flow**: Proto → SQL → Seeders  
✅ **Schema-aware**: Multi-schema support  
✅ **Checksum validation**: Prevents accidental changes  
✅ **Transactional**: Migrations rollback on error  
✅ **Migration ordering**: Via `migration_order` field  
✅ **Foreign key handling**: Two-phase creation  
✅ **Idempotent**: Safe to run multiple times  
✅ **Metadata tracking**: Full audit trail  
✅ **Rich annotations**: Comments, indexes, constraints  

---

## 📊 Proto Table Features

### Table-Level Options
- `table_name` - Explicit table name
- `schema_name` - Target schema
- `migration_order` - Creation order (default: 1000)
- `comment` - Table description
- `soft_delete` - Adds `deleted_at` column
- `audit_fields` - Adds `created_at`, `updated_at`
- `enable_rls` - Row-level security

### Column-Level Options
- `sql_type` - Override SQL type
- `primary_key` - Mark as primary key
- `not_null` - NOT NULL constraint
- `unique` - UNIQUE constraint
- `default_value` - Default value
- `comment` - Column description
- `foreign_key` - Foreign key relationship

### Index Options
- `name` - Index name
- `columns` - Indexed columns
- `unique` - Unique index
- `method` - BTREE, GIN, GIST, etc.
- `where` - Partial index condition

---

## 🐛 Troubleshooting

### Migration Fails

1. Check `public.schema_migrations` for error message
2. Fix the SQL file
3. Re-run migration (system tracks what succeeded)

### Schema Not Created

- Verify proto file has correct `schema_name` annotation
- Check proto package starts with `insuretech`
- Re-run migration system

### Foreign Key Errors

- Check `migration_order` values
- Referenced tables should have lower order numbers
- System creates FK in Phase 2 (after all tables exist)

---

## 📞 Support

For issues or questions:
- Check proto annotations in `proto/insuretech/common/v1/db.proto`
- Review existing proto entity files for examples
- Consult the migration logs in `schema_migrations` table

---

**Last Updated:** January 2025  
**Version:** 2.0 (Unified System)
