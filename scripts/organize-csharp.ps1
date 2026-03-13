# Organize C# generated files into folder structure based on namespaces
# Run after: buf generate

$csharpDir = Join-Path $PSScriptRoot "..\gen\csharp"

Write-Host "`nOrganizing C# files into folder structure..." -ForegroundColor Cyan

# Get all C# files
$files = Get-ChildItem -Path $csharpDir -Filter "*.cs" -File

$movedCount = 0

foreach ($file in $files) {
    # Read first few lines to find namespace
    $content = Get-Content $file.FullName -First 20
    $namespaceLine = $content | Where-Object { $_ -match "^namespace\s+(.+)" } | Select-Object -First 1
    
    if ($namespaceLine -match "namespace\s+(.+?)[\s;{]") {
        $namespace = $Matches[1].Trim()
        
        # Convert namespace to folder path
        # Example: Insuretech.Common.V1 -> Insuretech/Common/V1
        $folderPath = $namespace -replace '\.', '\'
        $targetDir = Join-Path $csharpDir $folderPath
        
        # Create directory if it doesn't exist
        if (-not (Test-Path $targetDir)) {
            New-Item -ItemType Directory -Path $targetDir -Force | Out-Null
        }
        
        # Move file
        $targetFile = Join-Path $targetDir $file.Name
        
        if ($file.FullName -ne $targetFile) {
            Move-Item -Path $file.FullName -Destination $targetFile -Force
            $movedCount++
            Write-Host "  ✓ $($file.Name) -> $folderPath\" -ForegroundColor Green
        }
    }
}

Write-Host "`nOrganized $movedCount C# files into folder structure." -ForegroundColor Green
Write-Host "Location: gen/csharp/" -ForegroundColor Cyan
