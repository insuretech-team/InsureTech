# DBManager — Complete CLI Guide

> **Database Management CLI** — dependency-aware primary/backup sync, health, and maintenance for PostgreSQL.

---

## Quick Start

```powershell
# Build
go build -o dbmanager.exe ./cmd/dbmanager/...

# Launch interactive TUI (no arguments)
./dbmanager.exe

# Show help
./dbmanager.exe --help

# Check status
./dbmanager.exe status

# Run migrations
./dbmanager.exe migrate --target=primary
```

### First-Time Setup

1. Create config file: `configs/database.yaml`
2. Set environment variables in `.env` at project root:
   ```env
   PGHOST=your-host
   PGPORT=5432
   PGDATABASE=your-db
   PGUSER=your-user
   PGPASSWORD=your-pass
   PGSSLMODE=require
   ```
3. Run: `./dbmanager.exe`

> DBManager auto-resolves `database.yaml` via the project-level config resolver and auto-loads `.env` from the project root.

---

## Invocation Modes

| Mode | Example | Notes |
|:-----|:--------|:------|
| **Interactive TUI** | `./dbmanager.exe` | No args → Bubble Tea command palette |
| **Cobra subcommand** | `./dbmanager.exe sync --commit` | Modern `--flag` style |
| **Legacy flag** | `./dbmanager.exe -cmd=sync -commit` | Backwards-compatible |

> **Tip:** Dashed pseudo-commands like `--migrate` are auto-rewritten to `migrate` with a hint.

---

## Global Flags

| Flag | Default | Description |
|:-----|:--------|:------------|
| `--config` | `database.yaml` | Path to database configuration file |
| `--source` | `primary` | Source database (`primary` / `backup`) |
| `--target` | `primary` | Target database (`primary` / `backup`) |

---

## Commands Reference

### Status & Diagnostics

#### `status`
Show database connection status and metrics for primary and backup.
```powershell
./dbmanager.exe status
```

#### `sizes`
Show approximate database and table sizes.
```powershell
./dbmanager.exe sizes
```

#### `test-init`
Verify database connections and query the `processes` table (legacy flags only).
```powershell
./dbmanager.exe -cmd=test-init
```

---

### Sync Operations

#### `sync`
Synchronize primary → backup in dependency-aware FK order with authoritative upsert.

| Flag | Default | Description |
|:-----|:--------|:------------|
| `--table` | *(all)* | Sync a single table only |
| `--commit` | `false` | Write changes (otherwise upsert-only) |
| `--prune` | `false` | Delete rows only in backup |
| `--fail-on-drift` | `false` | Exit non-zero if drift remains (CI) |
| `--report-format` | `table` | `table` / `markdown` / `csv` / `json` / `tui` |

```powershell
./dbmanager.exe sync                                         # dry-run
./dbmanager.exe sync --commit --prune                        # full sync
./dbmanager.exe sync --table=users --commit                  # single table
./dbmanager.exe sync --commit --prune --fail-on-drift        # CI mode
./dbmanager.exe sync --commit --report-format=json           # JSON report
./dbmanager.exe sync --commit --report-format=tui            # TUI report viewer
```

#### `sync-health-check`
Per-table row counts showing in-sync / out-of-sync status.
```powershell
./dbmanager.exe sync-health-check
```

#### `sync-repair`
Repair FK gaps for critical tables when sync fails due to FK issues.
```powershell
./dbmanager.exe sync-repair
```

#### `sync-users`
Synchronize user-related tables with special conflict resolution.
```powershell
./dbmanager.exe sync-users
```

---

### Migrations

#### `migrate`
Run proto-driven migrations with automatic pre-flight checks (proto freshness + SQL lint).

| Flag | Default | Description |
|:-----|:--------|:------------|
| `--target` | `primary` | `primary` / `backup` / `both` |
| `--prune` | `false` | Delete columns not in proto |
| `--strict` | `false` | Fail on schema drift (zombie columns, type mismatches) |

```powershell
./dbmanager.exe migrate --target=primary
./dbmanager.exe migrate --target=backup
./dbmanager.exe migrate --target=both
./dbmanager.exe migrate --target=primary --strict --prune
```

---

### Schema Inspection

#### `schema-discovery`
List all public base tables on primary.
```powershell
./dbmanager.exe schema-discovery
```

#### `schema-check`
Validate schema consistency between primary and backup. Suggests `rebuild-backup` if mismatched.
```powershell
./dbmanager.exe schema-check
```

#### `print-schema`
Print detailed schema information (tables, sizes, descriptions).

| Flag | Default | Description |
|:-----|:--------|:------------|
| `--schema` | *(all)* | Filter to specific schema |
| `--target` | `primary` | Database to inspect |

```powershell
./dbmanager.exe print-schema --target=primary
./dbmanager.exe print-schema --schema=auth --target=primary
```

<details>
<summary>Example output</summary>

```
╔═══════════════════════════════════════════════════════════════╗
║              SCHEMA INFORMATION - PRIMARY                     ║
╚═══════════════════════════════════════════════════════════════╝

📁 SCHEMA: auth
═══════════════════════════════════════════════════════════════

TABLE                          SIZE       DESCRIPTION
─────                          ────       ───────────
users                          8192 kB    Application users
customers                      512 kB     Customer entity
customer_addresses             256 kB     Addresses

Total tables: 15
```
</details>

#### `print-table`
Print comprehensive table info: columns, types, constraints, PKs, FKs, indexes, sizes. Supports `schema.table` format.

| Flag | Default | Description |
|:-----|:--------|:------------|
| `--table` | *(required)* | Table name (supports `schema.table`) |
| `--target` | `primary` | Database to inspect |

```powershell
./dbmanager.exe print-table --table=users --target=primary
./dbmanager.exe print-table --table=auth.users --target=primary
```

<details>
<summary>Example output</summary>

```
📋 Table: auth.users
   Description: Application users with authentication credentials

📊 Statistics:
   Rows: 1523  |  Table: 8.2 MB  |  Index: 2.1 MB  |  Total: 10.3 MB

🗂️  Columns (22):
NAME                TYPE              NULL  DEFAULT                  KEY
────                ────              ────  ───────                  ───
user_id             uuid              NO    uuid_generate_v4()       PK
email               text              NO                             UQ
password_hash       text              YES
...

🔗 Foreign Keys (0):

📇 Indexes (4):
   users_pkey (btree) on (user_id)
   users_email_key UNIQUE (btree) on (email)
   ...
```
</details>

#### `print-tables`
Print detailed info for **all** tables in a schema.

| Flag | Default | Description |
|:-----|:--------|:------------|
| `--schema` | *(all)* | Filter to a specific schema |
| `--target` | `primary` | Database to inspect |

```powershell
./dbmanager.exe print-tables --schema=auth --target=primary
./dbmanager.exe print-tables --target=primary              # all schemas
```

#### `print-all`
Print comprehensive database overview: all schemas, tables, sizes, row counts.
```powershell
./dbmanager.exe print-all --target=primary
./dbmanager.exe print-all --target=backup
```

#### `print-table-data`
Display actual table data in formatted tabular view.

| Flag | Default | Description |
|:-----|:--------|:------------|
| `--table` | *(required)* | Table name (supports `schema.table`) |
| `--target` | `primary` | Database to inspect |
| `--limit` | `100` | Max rows (max: 1000) |

```powershell
./dbmanager.exe print-table-data --table=auth.users --limit=20
./dbmanager.exe print-table-data --table=auth.customers --target=backup --limit=200
```

---

### SQL Execution

#### `sql`
Execute arbitrary SQL on primary, backup, or both.

| Flag | Default | Description |
|:-----|:--------|:------------|
| `--sql` | *(required)* | SQL query to execute |
| `--target` | `primary` | `primary` / `backup` / `both` |

```powershell
./dbmanager.exe sql --sql="SELECT 1" --target=primary
./dbmanager.exe sql --sql="DROP TABLE fabric_costs CASCADE;" --target=backup
./dbmanager.exe sql --sql="VACUUM ANALYZE;" --target=both
```

> ⚠️ No confirmation prompt for destructive statements — be careful.

---

### Failover & Recovery

| Command | Description |
|:--------|:------------|
| `failover` | Manual failover to backup database |
| `switchback` | Switch back to primary *(NYI)* |
| `rebuild-backup` | Rebuild backup schema to match primary |

```powershell
./dbmanager.exe failover
./dbmanager.exe rebuild-backup
```

---

### Backup & Restore

| Command | Description | Status |
|:--------|:------------|:-------|
| `backup` | Create compressed backup | *(NYI)* |
| `restore` | Restore from backup file | *(NYI)* |
| `list-backups` | List available backup files | ✅ |

```powershell
./dbmanager.exe list-backups
./dbmanager.exe restore --backup=backup_20240101.sql.gz --target=primary
```

---

### CSV Operations

#### `csv-backup` — Export tables to CSV
```powershell
./dbmanager.exe csv-backup --source=primary                    # all tables
./dbmanager.exe csv-backup --table=enquiries --source=primary  # single table
```

#### `csv-seed` — Import CSV into database
```powershell
./dbmanager.exe csv-seed --target=primary                      # all CSVs
./dbmanager.exe csv-seed --table=enquiries --target=primary    # single table
```

> ⚠️ `csv-seed` **truncates** the target table before importing.

CSV files are stored in `internal/db/backup/`.

---

### Data Copy & Comparison

#### `copy`
Copy data between databases.
```powershell
./dbmanager.exe copy --source=primary --target=backup
```

#### `compare`
Compare schemas in detail between primary and backup. *(NYI)*
```powershell
./dbmanager.exe compare
```

---

### Interactive Viewers

#### `view-tables`
Launch Bubble Tea interactive table viewer.
```powershell
./dbmanager.exe view-tables
./dbmanager.exe view-tables --table=enquiries
```

---

## Interactive TUI

Run `./dbmanager.exe` with **no arguments** to launch the interactive command palette.

### Features
- **Fuzzy filtering** — type to filter commands instantly
- **Background connection** — connects to databases automatically
- **Color-coded** — success/error states highlighted
- **All commands** — select and provide arguments inline

### Keyboard Shortcuts

| Key | Action |
|:----|:-------|
| `↑` / `↓` | Navigate menu |
| `Enter` | Select / execute command |
| Type | Filter commands |
| `Esc` | Return to menu from results |
| `Ctrl+C` / `q` | Quit |

### TUI Commands

| Command | Description |
|:--------|:------------|
| Status | Database connection status |
| Schema Discovery | List all tables |
| Schema Check | Validate consistency |
| Sync Health Check | Table sync status |
| Print All | All schemas and tables |
| Sizes | Database sizes |
| Exit | Quit |

> **Tip:** For commands needing extra flags (e.g., `--table`, `--sql`), use CLI mode instead.

### TUI Workflows

```
Check Health:     Launch → Status → Sync Health Check
Explore Schema:   Launch → Schema Discovery → Print All
Monitor Sizes:    Launch → Sizes
```

---

## Workflow Examples

### Daily Sync
```powershell
./dbmanager.exe sync-health-check         # 1. check status
./dbmanager.exe sync --commit --prune     # 2. full sync
./dbmanager.exe sync-health-check         # 3. verify
```

### Migration Pipeline
```powershell
./dbmanager.exe migrate --target=primary  # 1. migrate primary
./dbmanager.exe schema-check              # 2. verify schema
./dbmanager.exe rebuild-backup            # 3. fix backup if needed
./dbmanager.exe migrate --target=backup   # 4. migrate backup
./dbmanager.exe sync --commit --prune     # 5. sync data
```

### CI Pipeline
```powershell
./dbmanager.exe sync --commit --prune --fail-on-drift --report-format=json
```

### Schema Investigation
```powershell
./dbmanager.exe print-all --target=primary                           # overview
./dbmanager.exe print-table --table=auth.users --target=primary      # deep-dive
./dbmanager.exe print-table-data --table=auth.users --limit=20       # view data
```

### Schema Comparison (Primary vs Backup)
```powershell
./dbmanager.exe print-schema --schema=auth --target=primary > primary_auth.txt
./dbmanager.exe print-schema --schema=auth --target=backup  > backup_auth.txt
diff primary_auth.txt backup_auth.txt
```

### CSV Round-Trip
```powershell
./dbmanager.exe csv-backup --source=primary     # export
./dbmanager.exe csv-seed --target=backup         # import into backup
```

### Before/After Migration Verification
```powershell
./dbmanager.exe print-table --table=auth.customers  # before
./dbmanager.exe migrate --target=primary
./dbmanager.exe print-table --table=auth.customers  # after
```

---

## Pro Tips

1. **Use `schema.table` format** for clarity: `--table=auth.users`
2. **Redirect to file** for documentation: `> output.txt`
3. **Pipe to grep** for filtering: `| grep -i "index"`
4. **Shell aliases** for quick access:
   ```powershell
   Set-Alias dbm './dbmanager.exe'
   function dbm-status { ./dbmanager.exe status }
   function dbm-sync { ./dbmanager.exe sync --commit }
   ```

---

## Report Formats (`--report-format`)

| Format | Description |
|:-------|:------------|
| `table` | *(default)* Aligned columns via tabwriter |
| `markdown` | GitHub-flavored markdown table |
| `csv` | Comma-separated values |
| `json` | Full metadata: started, duration, tables, pruned |
| `tui` | Interactive Bubble Tea viewer |

---

## Troubleshooting

### Connection Issues
| Problem | Solution |
|:--------|:---------|
| TUI doesn't start | Verify exe exists, rebuild with `go build -o dbmanager.exe ./cmd/dbmanager/...` |
| Can't connect to database | Check `.env` in project root and `configs/database.yaml` |
| "Database manager not initialized" | Verify database configuration in `database.yaml` |
| Network errors | Check firewall, VPN, and database server accessibility |

### Data Issues
| Problem | Solution |
|:--------|:---------|
| "Table does not exist" | Use `print-schema` to list available tables |
| "Cannot count rows" | Permission issues; row count shows "N/A" |
| Slow performance | Use `--schema` flag to filter scope |
| Commands show no output | Check result view for success/error messages |

---

## NYI Commands

| Command | Description |
|:--------|:------------|
| `backup` | Create compressed database backup |
| `restore` | Restore from backup file |
| `compare` | Detailed schema diff between DBs |
| `switchback` | Switch back to primary after failover |
