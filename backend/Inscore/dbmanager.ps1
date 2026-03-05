#!/usr/bin/env pwsh
# DBManager wrapper - Run from backend/inscore directory
# Usage: ./dbmanager.ps1 <command> [args]
# Example: ./dbmanager.ps1 migrate --target=primary

$args_string = $args -join ' '
go run ./cmd/dbmanager $args_string
