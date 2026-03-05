# =====================================================
# Safe Database Reset (Primary/Backup)
# =====================================================
# Deletes all tables in public schema, then drops all other schemas.
# Keeps only empty public schema and system schemas.
# =====================================================

param(
    [Parameter(Mandatory=$false)]
    [ValidateSet("primary","backup","both")]
    [string]$Target = "both",

    [switch]$DryRun,
    [switch]$SkipConfirmation
)

$ErrorActionPreference = "Stop"

# SQL to drop all tables in public schema, then drop all custom schemas
$nukeSQL = @"
DO `$`$
DECLARE
    r RECORD;
BEGIN
    -- Drop all tables in public schema
    FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = 'public') LOOP
        EXECUTE 'DROP TABLE IF EXISTS public.' || quote_ident(r.tablename) || ' CASCADE';
    END LOOP;
    
    -- Drop all non-system schemas (keeping public, pg_catalog, information_schema, pg_toast)
    FOR r IN (
        SELECT schema_name 
        FROM information_schema.schemata 
        WHERE schema_name NOT IN ('public', 'pg_catalog', 'information_schema', 'pg_toast')
        AND schema_name NOT LIKE 'pg_%'
    ) LOOP
        EXECUTE 'DROP SCHEMA IF EXISTS ' || quote_ident(r.schema_name) || ' CASCADE';
    END LOOP;
END
`$`$;
"@

Write-Host "" 
Write-Host "========================================" -ForegroundColor Yellow
Write-Host "SAFE RESET (DROP CUSTOM SCHEMAS ONLY)" -ForegroundColor Yellow
Write-Host "========================================" -ForegroundColor Yellow
Write-Host "Target: $Target" -ForegroundColor Cyan
Write-Host "DryRun: $DryRun" -ForegroundColor Cyan
Write-Host "" 

if (-not $SkipConfirmation -and -not $DryRun) {
  $confirmation = Read-Host "Type 'RESET' to confirm"
  if ($confirmation -ne "RESET") {
    Write-Host "Aborted." -ForegroundColor Red
    exit 0
  }
}

if ($DryRun) {
  Write-Host "--- SQL Preview ---" -ForegroundColor Gray
  Write-Host $nukeSQL -ForegroundColor Gray
  exit 0
}

$targets = if ($Target -eq "both") { @("primary","backup") } else { @($Target) }

$root = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location "$root\backend\inscore"

foreach ($db in $targets) {
  Write-Host "Resetting $db..." -ForegroundColor Yellow
  
  # Save SQL to temp file to avoid escaping issues
  $tempSqlFile = [System.IO.Path]::GetTempFileName() + ".sql"
  $nukeSQL | Out-File -FilePath $tempSqlFile -Encoding UTF8 -NoNewline
  
  try {
    # Read the SQL from file and execute
    $sqlContent = Get-Content $tempSqlFile -Raw
    go run ./cmd/dbmanager sql --sql=$sqlContent --target=$db
    
    if ($LASTEXITCODE -ne 0) { 
      Write-Host "Error: Reset failed for $db" -ForegroundColor Red
      Write-Host "Exit code: $LASTEXITCODE" -ForegroundColor Red
      throw "Reset failed for $db" 
    }
    
    Write-Host "✓ Reset complete for $db" -ForegroundColor Green
  }
  finally {
    # Clean up temp file
    if (Test-Path $tempSqlFile) {
      Remove-Item $tempSqlFile -Force
    }
  }
}
