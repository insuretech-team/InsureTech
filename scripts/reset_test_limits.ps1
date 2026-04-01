# ============================================================
# reset_test_limits.ps1 -- Reset OTP/login rate limits for testing
#
# Usage:
#   .\scripts\reset_test_limits.ps1
#   .\scripts\reset_test_limits.ps1 -MobileNumber "+8801712345678"
#   .\scripts\reset_test_limits.ps1 -MobileNumber "+8801712345678" -All
#   .\scripts\reset_test_limits.ps1 -DryRun
# ============================================================
param(
    [string]$MobileNumber = "+8801347201751",
    [switch]$All,
    [switch]$DryRun
)

$Remote    = "insureadmin@146.190.97.242"
$SshBase   = "ssh -o StrictHostKeyChecking=no $Remote"
$DbUser    = "insuretech"
$DbName    = "insuretech"

# Normalize mobile: strip leading +
$MobileStripped = $MobileNumber.TrimStart('+')
$MobileFull     = if ($MobileNumber.StartsWith('+')) { $MobileNumber } else { "+$MobileNumber" }

Write-Host ""
Write-Host "================================================================" -ForegroundColor Cyan
Write-Host "  InsureTech -- Rate Limit Reset Tool (Dev/Test Only)" -ForegroundColor Cyan
Write-Host "================================================================" -ForegroundColor Cyan
Write-Host "  Mobile   : $MobileFull"
Write-Host "  All      : $($All.IsPresent)"
Write-Host "  Dry run  : $($DryRun.IsPresent)"
Write-Host ""

# Build Redis keys to clear
$otpTypes = @("login", "registration", "reset_password")
$redisKeys = @()
foreach ($t in $otpTypes) {
    $redisKeys += "otp_rl:min:${t}:${MobileStripped}"
    $redisKeys += "otp_rl:min:${t}:${MobileFull}"
    $redisKeys += "otp_rl:hour:${t}:${MobileStripped}"
    $redisKeys += "otp_rl:hour:${t}:${MobileFull}"
    $redisKeys += "otp_rl:day:${t}:${MobileStripped}"
    $redisKeys += "otp_rl:day:${t}:${MobileFull}"
}
$redisKeys += "otp_cooldown:${MobileStripped}"
$redisKeys += "otp_cooldown:${MobileFull}"

if ($DryRun) {
    Write-Host "[DRY RUN] Redis keys that would be deleted:" -ForegroundColor Yellow
    foreach ($k in $redisKeys) { Write-Host "  DEL $k" }
    Write-Host ""
    Write-Host "[DRY RUN] DB: UPDATE authn_schema.otps SET expires_at=NOW()-2h WHERE recipient IN ('$MobileStripped','$MobileFull')" -ForegroundColor Yellow
    if ($All) {
        Write-Host "[DRY RUN] DB: DELETE FROM authn_schema.login_attempts WHERE mobile_number IN (...) AND created_at > NOW()-1h" -ForegroundColor Yellow
    }
    Write-Host ""
    Write-Host "[DRY RUN] No changes made." -ForegroundColor Yellow
    exit 0
}

# Step 1: Clear Redis OTP rate limit keys
Write-Host ">> [1/3] Clearing Redis OTP rate limit keys..." -ForegroundColor Cyan
$cleared = 0
foreach ($key in $redisKeys) {
    $cmd = "docker exec insuretech-redis redis-cli DEL $key"
    $result = wsl -d Ubuntu bash -c "ssh -o StrictHostKeyChecking=no $Remote '$cmd'"
    if ($result -match "^1$") {
        Write-Host "    DEL $key -> removed" -ForegroundColor Green
        $cleared++
    }
}
if ($cleared -eq 0) {
    Write-Host "    No Redis OTP rate limit keys found (already clear)" -ForegroundColor Gray
} else {
    Write-Host "    Cleared $cleared Redis key(s)" -ForegroundColor Green
}

# Step 2: Expire OTP records in DB
# Uses psql on the remote server with the DATABASE_URL from .env
Write-Host ">> [2/3] Expiring OTP records in DB..." -ForegroundColor Cyan

$RemoteDir = "/home/insureadmin/insuretech"
$sqlOtp = "UPDATE authn_schema.otps SET expires_at = NOW() - INTERVAL '2 hours' WHERE recipient IN ('$MobileStripped', '$MobileFull') AND expires_at > NOW();"

# Write bash script with LF line endings (Unix) to avoid \r issues
$otpLines = @(
    "#!/bin/bash",
    "cd $RemoteDir",
    "set -a; source .env; set +a",
    "echo `"$sqlOtp`" | psql `"`$DATABASE_URL`""
)
$otpScriptContent = $otpLines -join "`n"
$tmpOtp = [System.IO.Path]::GetTempFileName()
[System.IO.File]::WriteAllText($tmpOtp, $otpScriptContent)

& scp.exe -o StrictHostKeyChecking=no $tmpOtp "${Remote}:/tmp/reset_otp.sh"
$dbResult = & ssh.exe -o StrictHostKeyChecking=no $Remote "bash /tmp/reset_otp.sh; rm -f /tmp/reset_otp.sh" 2>&1
Remove-Item $tmpOtp -ErrorAction SilentlyContinue

if ($dbResult -match "UPDATE") {
    Write-Host "    OTPs expired in DB: $dbResult" -ForegroundColor Green
} else {
    Write-Host "    DB result: $dbResult" -ForegroundColor Yellow
}

# Step 3: Clear login attempts (optional -All)
if ($All) {
    Write-Host ">> [3/3] Clearing recent login attempts in DB..." -ForegroundColor Cyan
    $sqlLogin = "DELETE FROM authn_schema.login_attempts WHERE mobile_number IN ('$MobileStripped', '$MobileFull') AND created_at > NOW() - INTERVAL '1 hour';"

    $loginLines = @(
        "#!/bin/bash",
        "cd $RemoteDir",
        "set -a; source .env; set +a",
        "echo `"$sqlLogin`" | psql `"`$DATABASE_URL`""
    )
    $loginScriptContent = $loginLines -join "`n"
    $tmpLogin = [System.IO.Path]::GetTempFileName()
    [System.IO.File]::WriteAllText($tmpLogin, $loginScriptContent)

    & scp.exe -o StrictHostKeyChecking=no $tmpLogin "${Remote}:/tmp/reset_login.sh"
    $loginResult = & ssh.exe -o StrictHostKeyChecking=no $Remote "bash /tmp/reset_login.sh; rm -f /tmp/reset_login.sh" 2>&1
    Remove-Item $tmpLogin -ErrorAction SilentlyContinue

    if ($loginResult -match "DELETE") {
        Write-Host "    Login attempts cleared: $loginResult" -ForegroundColor Green
    } else {
        Write-Host "    DB result: $loginResult" -ForegroundColor Yellow
    }
} else {
    Write-Host ">> [3/3] Skipping login attempt reset (use -All to include)" -ForegroundColor Gray
}

Write-Host ""
Write-Host "================================================================" -ForegroundColor Green
Write-Host "  Done! Rate limits cleared for: $MobileFull" -ForegroundColor Green
Write-Host "================================================================" -ForegroundColor Green
Write-Host ""
Write-Host "You can now run OTP send + login tests again." -ForegroundColor Cyan
Write-Host ""
