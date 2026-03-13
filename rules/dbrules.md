# Database Rules

> Quick reference for database design and migration standards.

---

## Naming Conventions

| Element | Pattern | Example |
|---------|---------|---------|
| **Primary Key** | `{table_singular}_id` | `tenant_id`, `user_id`, `policy_id` |
| **Foreign Key** | `{referenced_table}_id` | `customer_id`, `created_by` |
| **Table Name** | `snake_case`, plural | `users`, `policy_riders`, `fraud_alerts` |
| **Column Name** | `snake_case` | `created_at`, `login_attempts` |
| **Index** | `idx_{table}_{columns}` | `idx_policies_customer_status` |
| **Constraint** | `chk_{table}_{column}` | `chk_users_email` |

---

## Core Rules

| # | Rule | Description |
|---|------|-------------|
| 1 | **Single Source of Truth** | All schema definitions in Proto, never manual DDL |
| 2 | **Consistent ID Naming** | Use `{entity}_id` not generic `id` (e.g., `tenant_id`) |
| 3 | **No Hardcoded DDL** | SQL migrations = data only, no `CREATE TABLE` |
| 4 | **No Constraint Drift** | `ON DELETE`, `NOT NULL`, `CHECK` must match Proto exactly |
| 5 | **No Type Mismatch** | Column types must match Proto-derived types |
| 6 | **Zombie Detection** | Columns in DB but not Proto are flagged/pruned |
| 7 | **Code Freshness** | Block migration if `*.pb.go` older than `*.proto` |
| 8 | **Auto Type Sync** | Engine auto-alters columns to match Proto |
| 9 | **CI Schema Diff** | Fail deployment on unexpected schema drift |
| 10 | **AuditInfo as JSONB** | All entities store `audit_info` as JSONB column |
| 11 | **Money + Currency** | Every Money field has companion `_currency` column |
| 12 | **CHECK Constraints** | Use `check_constraint` option in Proto for validation |

---

## Enforcement

| Rule | Action | Flag |
|------|--------|------|
| Zombie Columns | Warn → Prune → Fail | `--prune`, `--strict` |
| Type Mismatch | Auto-heal via ALTER | Default |
| Constraint Drift | Auto-repair FK rules | Default |
| Code Freshness | Block execution | Default |
| Schema Drift | CI failure | `--strict` |

---

## ID Naming Migration

**Bad** → **Good**:
```diff
- id UUID PRIMARY KEY
+ tenant_id UUID PRIMARY KEY

- id UUID PRIMARY KEY  
+ fraud_rule_id UUID PRIMARY KEY

- id UUID PRIMARY KEY
+ workflow_definition_id UUID PRIMARY KEY
```

---

## Quick Checklist

- [ ] Every table has `{entity}_id` as PK
- [ ] Every FK references explicit `{table}_id`
- [ ] Money fields have `_currency` companion
- [ ] AuditInfo has JSONB column annotation
- [ ] CHECK constraints defined in Proto
- [ ] No zombie columns in production