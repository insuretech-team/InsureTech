# `dbx` — InsureTech Database CLI

> **dbx** *(Database eXplorer)* — dependency-aware sync, live inspection, migrations, and interactive TUI for InsureTech PostgreSQL databases.

---

## Quick Start

```bash
# Build
go build -o dbx ./backend/inscore/cmd/dbx/...

# Alias for convenience (add to your shell profile)
alias dbx='go run ./backend/inscore/cmd/dbx'

# Launch interactive TUI
dbx

# Check status
dbx status

# Run migrations
dbx migrate --target=primary
```

### First-Time Setup

1. Copy `.env.example` → `.env` at project root and fill in:
   ```env
   PGHOST=your-host
   PGPORT=5432
   PGDATABASE=your-db
   PGUSER=your-user
   PGPASSWORD=your-pass
   PGSSLMODE=require
   ```
2. Ensure `configs/database.yaml` exists in the inscore backend.
3. Run `dbx` — it auto-resolves config and loads `.env` from the project root.

---

## Invocation Modes

| Mode | Command | Notes |
|:-----|:--------|:------|
| **Interactive TUI** | `dbx` | No args → guided command palette with dropdowns |
| **Cobra subcommand** | `dbx sync --commit` | Modern `--flag` style |
| **Legacy flag** | `dbx -cmd=sync -commit` | Backwards-compatible |

> **Tip:** Dashed pseudo-commands like `--migrate` are auto-rewritten to `migrate`.

---

## Global Flags

| Flag | Default | Description |
|:-----|:--------|:------------|
| `--config` | `database.yaml` | Path to database configuration file |
| `--source` | `primary` | Source database (`primary` / `backup`) |
| `--target` | `primary` | Target database (`primary` / `backup`) |

---

## Command Reference

### 🟢 Status & Diagnostics

#### `status`
Show connection status and metrics for primary and backup databases.
```bash
dbx status
```
**Output:** connectivity state, connection counts, failover count, last sync/backup times.

---

#### `sizes`
Show approximate database and table sizes for capacity planning.
```bash
dbx sizes
```

---

#### `sync-health-check`
Per-table row counts comparing primary vs backup — shows in-sync / out-of-sync.
```bash
dbx sync-health-check
```
**Use this before and after a sync** to verify what changed.

---

#### `schema-discovery`
List all public base tables on the primary database.
```bash
dbx schema-discovery
```

---

#### `schema-check`
Validate schema consistency between primary and backup. Suggests `rebuild-backup` if mismatched.
```bash
dbx schema-check
```

---

### 🔄 Sync Operations

#### `sync`
Synchronize primary → backup in dependency-aware FK order with authoritative upsert.

| Flag | Default | Description |
|:-----|:--------|:------------|
| `--table` | *(all tables)* | Sync a single table only |
| `--commit` | `false` | Write changes (omit for dry-run) |
| `--prune` | `false` | Delete rows that only exist in backup |
| `--fail-on-drift` | `false` | Exit non-zero if drift remains (for CI) |
| `--report-format` | `table` | `table` / `markdown` / `csv` / `json` / `tui` |

```bash
dbx sync                                          # dry-run — shows what would change
dbx sync --commit                                 # upsert only
dbx sync --commit --prune                         # full authoritative sync
dbx sync --table=auth.users --commit              # single table sync
dbx sync --commit --prune --fail-on-drift         # CI mode — fails if drift remains
dbx sync --commit --report-format=json            # JSON report for automation
dbx sync --commit --report-format=tui             # interactive TUI report viewer
```

---

#### `sync-repair`
Repair FK gaps for critical tables when sync fails due to foreign key violations.
```bash
dbx sync-repair
```

---

#### `sync-users`
Synchronize user-related tables with special conflict resolution for unique constraints.
```bash
dbx sync-users
```

---

### 🚀 Migrations

#### `migrate`
Run proto-driven SQL migrations with automatic pre-flight checks (proto freshness + SQL lint).

| Flag | Default | Description |
|:-----|:--------|:------------|
| `--target` | `primary` | `primary` / `backup` / `both` |
| `--prune` | `false` | Drop columns not present in proto definition |
| `--strict` | `false` | Fail on schema drift (zombie columns, type mismatches) |

```bash
dbx migrate --target=primary                       # standard migration
dbx migrate --target=backup                        # migrate backup only
dbx migrate --target=both                          # migrate both databases
dbx migrate --target=primary --strict              # fail on any drift
dbx migrate --target=primary --strict --prune      # strict + remove zombie columns
```

> **Migration pipeline:** `migrate primary` → `schema-check` → `rebuild-backup` → `migrate backup` → `sync`

---

### 🔍 Schema Inspection

#### `print-schema`
Print detailed schema information: tables, sizes, descriptions.

| Flag | Default | Description |
|:-----|:--------|:------------|
| `--schema` | *(all schemas)* | Filter to a specific schema |
| `--target` | `primary` | Database to inspect |

```bash
dbx print-schema --target=primary                  # all schemas
dbx print-schema --schema=public --target=primary  # public schema only
dbx print-schema --schema=auth --target=backup     # backup auth schema
```

**TUI:** Select `Print Schema` → pick schema from live dropdown → pick target

---

#### `print-tables`
Print detailed info for **all** tables in a schema (columns, constraints, FKs, indexes, sizes).

| Flag | Default | Description |
|:-----|:--------|:------------|
| `--schema` | *(all schemas)* | Filter to a specific schema |
| `--target` | `primary` | Database to inspect |

```bash
dbx print-tables --schema=public --target=primary  # all tables in public
dbx print-tables --target=primary                  # all schemas, all tables
```

**TUI:** Select `Print Tables` → pick schema → pick target

---

#### `print-table`
Print comprehensive info for a single table: columns, types, constraints, PKs, FKs, indexes, row count, size.

Supports `schema.table` format.

| Flag | Default | Description |
|:-----|:--------|:------------|
| `--table` | *(required)* | Table name — supports `schema.table` |
| `--target` | `primary` | Database to inspect |

```bash
dbx print-table --table=users --target=primary
dbx print-table --table=public.users --target=primary
dbx print-table --table=auth.customers --target=backup
```

**TUI:** Select `Print Table` → pick schema → table list auto-populates → pick table → pick target
→ produces `print-table --table=auth.users --target=primary`

<details>
<summary>Example output</summary>

```
📋 Table: auth.users
   Description: Application users

📊 Statistics:
   Rows: 1523  |  Table: 8.2 MB  |  Index: 2.1 MB  |  Total: 10.3 MB

🗂️  Columns (22):
NAME                TYPE              NULL  DEFAULT                  KEY
────                ────              ────  ───────                  ───
user_id             uuid              NO    uuid_generate_v4()       PK
email               text              NO                             UQ
password_hash       text              YES

🔗 Foreign Keys (2):
   tenant_id → tenants.tenant_id

📇 Indexes (4):
   users_pkey (btree) on (user_id)
   users_email_key UNIQUE (btree) on (email)
```
</details>

---

#### `print-all`
Print comprehensive database overview: all schemas, tables, sizes, row counts.

| Flag | Default | Description |
|:-----|:--------|:------------|
| `--target` | `primary` | Database to inspect |

```bash
dbx print-all --target=primary
dbx print-all --target=backup
```

**TUI:** Select `Print All` → pick target

---

#### `print-table-data`
Display actual rows from a table in a formatted tabular view.

| Flag | Default | Description |
|:-----|:--------|:------------|
| `--table` | *(required)* | Table name — supports `schema.table` |
| `--target` | `primary` | Database to inspect |
| `--limit` | `100` | Max rows to display (max: 1000) |

```bash
dbx print-table-data --table=auth.users --limit=20
dbx print-table-data --table=public.products --target=backup --limit=200
```

**TUI:** Select `Print Table Data` → pick schema → table auto-populates → pick table → enter limit → pick target
→ produces `print-table-data --table=auth.users --limit=20 --target=primary`

---

### 💾 SQL Execution

#### `sql`
Execute arbitrary SQL on primary, backup, or both.

| Flag | Default | Description |
|:-----|:--------|:------------|
| `--sql` | *(required)* | SQL query to execute |
| `--target` | `primary` | `primary` / `backup` / `both` |

```bash
dbx sql --sql="SELECT COUNT(*) FROM auth.users" --target=primary
dbx sql --sql="SELECT version()" --target=both
dbx sql --sql="VACUUM ANALYZE;" --target=primary
dbx sql --sql="DROP TABLE fabric_costs CASCADE;" --target=backup
```

**TUI:** Select `SQL Query` → pick target → type your query

> ⚠️ No confirmation prompt for destructive statements. Use `--target=backup` for safety testing.

---

### 🛡️ Failover & Recovery

| Command | Description | Status |
|:--------|:------------|:-------|
| `failover` | Manual switch to backup database | ✅ |
| `switchback` | Switch back to primary | *(NYI)* |
| `rebuild-backup` | Rebuild backup schema to match primary | ✅ |

```bash
dbx failover                   # switch to backup (emergency use)
dbx rebuild-backup             # re-align backup schema after drift
```

---

### 📦 Backup & Restore

| Command | Description | Status |
|:--------|:------------|:-------|
| `backup` | Create compressed database backup | *(NYI)* |
| `restore` | Restore from backup file | *(NYI)* |
| `list-backups` | List available backup files | ✅ |

```bash
dbx list-backups
dbx restore --backup=backup_20240101.sql.gz --target=primary
```

---

### 📋 CSV Operations

#### `csv-backup` — Export tables to CSV
```bash
dbx csv-backup --source=primary                    # export all tables
dbx csv-backup --table=auth.users --source=primary # export single table
```

#### `csv-seed` — Import CSV into database
```bash
dbx csv-seed --target=primary                      # import all CSVs
dbx csv-seed --table=auth.users --target=primary   # import single table
```

> ⚠️ `csv-seed` **truncates** the target table before importing. CSV files are stored in `internal/db/backup/`.

---

### 📊 Data Copy & Comparison

#### `copy`
Copy data between databases.
```bash
dbx copy --source=primary --target=backup
```

#### `compare`
Detailed schema diff between primary and backup. *(NYI)*
```bash
dbx compare
```

---

### 🖥️ Interactive Viewers

#### `view-tables`
Launch Bubble Tea interactive table viewer with CSV export integration.
```bash
dbx view-tables
dbx view-tables --table=auth.users
```

---

## 🎮 Interactive TUI

Run `dbx` with **no arguments** to launch the full interactive TUI.

```
┌─────────────────────────────────────────────┐
│  DBManager Interactive CLI                  │
│  Connected to database          ✓           │
│                                             │
│  ┌─────────────────────────────────────┐   │
│  │ Type command or '/' for menu        │   │
│  └─────────────────────────────────────┘   │
│                                             │
│  '/' = menu  •  'help' = reference         │
└─────────────────────────────────────────────┘
```

### TUI FSM States

```
Not connected ──/──► Connecting ──► Menu
                                     │
                              Select command
                                     │
                    ┌────────────────┴───────────────────┐
                    │ Has args?                          │
                    ▼ YES                        NO ▼
                  Form                        Executing
                 (dropdowns)                    │
                    │                           │
                    └──────────► Pager ◄────────┘
                              (scrollable result)
                                     │
                                    Esc
                                     │
                                   Input
```

### Keyboard Shortcuts

| Key | Screen | Action |
|:----|:-------|:-------|
| `/` | Input | Open command menu |
| `↑` / `↓` | Menu / Form | Navigate list |
| `Enter` | Any | Select / confirm / execute |
| `Esc` | Menu | Back to input |
| `Esc` | Form | Back one step (or back to menu) |
| `Esc` / `q` | Pager | Back to input |
| `j` / `k` | Pager | Scroll up / down |
| `PgUp` / `PgDn` | Pager | Page scroll |
| `Ctrl+C` | Any | Quit |
| `help` | Input | Show full command reference |

### Commands Available in TUI

| Menu Item | Form Steps | Produces |
|:----------|:-----------|:---------|
| Status | *(none — executes immediately)* | `status` |
| Schema Discovery | *(none)* | `schema-discovery` |
| Schema Check | *(none)* | `schema-check` |
| Sync Health Check | *(none)* | `sync-health-check` |
| Sizes | *(none)* | `sizes` |
| Print Schema | Schema (live) → Target | `print-schema --schema=X --target=Y` |
| Print Tables | Schema (live) → Target | `print-tables --schema=X --target=Y` |
| Print Table | Schema (live) → Table (live cascade) → Target | `print-table --table=schema.table --target=Y` |
| Print All | Target | `print-all --target=Y` |
| Print Table Data | Schema → Table (cascade) → Limit → Target | `print-table-data --table=schema.table --limit=N --target=Y` |
| SQL Query | Target → SQL text | `sql --sql=... --target=Y` |
| Help | *(none)* | Shows full reference in pager |
| Exit | *(none)* | Quits |

> **Live dropdowns:** Schema and table lists are loaded from the database after connecting. Tables auto-filter by the schema you pick — you always get accurate, up-to-date options.

> **CLI-only commands:** `sync`, `migrate`, `csv-backup`, `csv-seed`, `failover`, `rebuild-backup` are not available in TUI interactive mode — use CLI directly for these.

---

## Workflow Examples

### Daily Health Check
```bash
dbx status                     # connection status
dbx sync-health-check          # row count comparison
```

### Full Sync Pipeline
```bash
dbx sync-health-check          # 1. baseline check
dbx sync --commit --prune      # 2. full authoritative sync
dbx sync-health-check          # 3. verify in sync
```

### Migration Pipeline
```bash
dbx migrate --target=primary   # 1. migrate primary
dbx schema-check               # 2. verify schema match
dbx rebuild-backup             # 3. fix backup schema if needed
dbx migrate --target=backup    # 4. migrate backup
dbx sync --commit --prune      # 5. sync data
```

### Schema Investigation
```bash
dbx print-all --target=primary                          # full overview
dbx print-schema --schema=auth --target=primary         # schema detail
dbx print-table --table=auth.users --target=primary     # table deep-dive
dbx print-table-data --table=auth.users --limit=20      # view actual rows
```

### Primary vs Backup Comparison
```bash
dbx print-schema --schema=auth --target=primary > primary_auth.txt
dbx print-schema --schema=auth --target=backup  > backup_auth.txt
diff primary_auth.txt backup_auth.txt
```

### CSV Round-Trip
```bash
dbx csv-backup --source=primary       # export all tables to CSV
dbx csv-seed --target=backup          # import into backup
```

### CI Pipeline
```bash
dbx sync --commit --prune --fail-on-drift --report-format=json
```

### Before/After Migration Verification
```bash
dbx print-table --table=auth.customers  # snapshot before
dbx migrate --target=primary
dbx print-table --table=auth.customers  # verify after
```

---

## Shell Aliases (Recommended)

**bash / zsh** (add to `~/.bashrc` or `~/.zshrc`):
```bash
alias dbx='go run ./backend/inscore/cmd/dbx'
alias dbx-status='dbx status'
alias dbx-sync='dbx sync --commit --prune'
alias dbx-health='dbx sync-health-check'
alias dbx-migrate='dbx migrate --target=primary'
```

**PowerShell** (add to `$PROFILE`):
```powershell
function dbx { go run ./backend/inscore/cmd/dbx @args }
function dbx-status  { dbx status }
function dbx-sync    { dbx sync --commit --prune }
function dbx-health  { dbx sync-health-check }
function dbx-migrate { dbx migrate --target=primary }
```

---

## Report Formats

Available for `sync --report-format`:

| Format | Description |
|:-------|:------------|
| `table` | *(default)* Aligned columns |
| `markdown` | GitHub-flavored markdown table |
| `csv` | Comma-separated for spreadsheets |
| `json` | Full metadata: started, duration, tables, pruned |
| `tui` | Interactive Bubble Tea viewer |

---

## Troubleshooting

| Problem | Solution |
|:--------|:---------|
| Can't connect | Check `.env` at project root and `configs/database.yaml` |
| "no required module provides package" | Run from project root or ensure `go.work` is present |
| "Database manager not initialized" | Verify `database.yaml` and DB credentials |
| Table dropdown empty in TUI | DB connected but schema has no base tables — check schema name |
| Schema list shows only "public" | Cache load failed — check DB permissions on `information_schema` |
| Slow `print-tables` | Use `--schema=X` to filter scope |
| `csv-seed` data missing | Check CSV file exists in `internal/db/backup/` |

---

## NYI Commands

| Command | Description |
|:--------|:------------|
| `backup` | Create compressed database backup |
| `restore` | Restore from backup file |
| `compare` | Detailed schema diff between DBs |
| `switchback` | Switch back to primary after failover |
