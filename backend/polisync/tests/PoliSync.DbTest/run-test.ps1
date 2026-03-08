# Load environment variables from .env file and run database test
Write-Host "=== PoliSync Database Live Test ===" -ForegroundColor Cyan
Write-Host "Loading environment variables from .env file...`n" -ForegroundColor Yellow

# Load .env file
$envFile = "../../../../.env"
if (Test-Path $envFile) {
    Get-Content $envFile | ForEach-Object {
        if ($_ -match '^\s*([^#][^=]+)=(.*)$') {
            $key = $matches[1].Trim()
            $value = $matches[2].Trim().Trim('"').Trim("'")
            [Environment]::SetEnvironmentVariable($key, $value, "Process")
        }
    }
    Write-Host "✅ Environment variables loaded" -ForegroundColor Green
} else {
    Write-Host "❌ .env file not found at: $envFile" -ForegroundColor Red
    exit 1
}

# Display connection info (without password)
Write-Host "`nDatabase Connection:" -ForegroundColor Cyan
Write-Host "  Host: $env:PGHOST"
Write-Host "  Port: $env:PGPORT"
Write-Host "  Database: $env:PGDATABASE"
Write-Host "  User: $env:PGUSER"
Write-Host "  SSL Mode: $env:PGSSLMODE"
Write-Host ""

# Build and run the test
Write-Host "Building test project..." -ForegroundColor Yellow
dotnet build --configuration Release

if ($LASTEXITCODE -eq 0) {
    Write-Host "`nRunning database tests...`n" -ForegroundColor Yellow
    dotnet run --configuration Release --no-build
} else {
    Write-Host "`n❌ Build failed!" -ForegroundColor Red
    exit 1
}
