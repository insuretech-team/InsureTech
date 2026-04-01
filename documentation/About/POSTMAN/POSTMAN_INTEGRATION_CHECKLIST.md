# InsureTech Postman Integration Checklist

## Executive Summary

This document provides a step-by-step implementation guide for integrating advanced Postman collection generation and testing into the InsureTech project's CI/CD pipelines. The implementation enables:

- ✅ Automatic generation of 1000+ API endpoints with test scripts
- ✅ Intelligent environment variable propagation (tokens auto-captured)
- ✅ Support for both B2C JWT and B2B server-side session flows
- ✅ Seamless Postman API integration with `POSTMAN_API_KEY`
- ✅ Newman CLI support for headless testing
- ✅ Full pipeline integration (bash + PowerShell)

## Current Status

### Already Implemented
- ✅ `sync_postman.py` — Generates collections and environment files
- ✅ `run_api_pipeline.sh` — Step 15 includes Postman generation
- ✅ Environment file structure with comprehensive variables
- ✅ Auth smoke tests (`auth_smoke.postman_collection.json`)
- ✅ Authorization enforcement tests (`b2c_authz_enforcement.postman_collection.json`)

### Ready for Implementation
- 📋 Enhanced environment propagation scripts (provided in `sync_postman_enhanced.py`)
- 📋 PowerShell integration patch (`POSTMAN_MIGRATION_PATCH.ps1`)
- 📋 Complete user documentation (`POSTMAN_COLLECTION_GUIDE.md`)
- 📋 CI/CD integration examples

## Implementation Phases

### Phase 1: Preparation (5-10 minutes)

#### Step 1.1: Review Current State
```bash
cd E:\Projects\InsureTech

# Check existing Postman files
ls -la api/postman/

# Verify sync_postman.py exists
ls -la api/generator/sync_postman.py

# Check pipeline scripts
ls -la run_api_pipeline.sh run_migration.ps1
```

**Expected Output:**
```
✓ api/postman/ contains:
  - InsureTech.postman_collection.json (if generated)
  - InsureTech_local.postman_environment.json
  - auth_smoke.postman_collection.json
  - b2c_authz_enforcement.postman_collection.json

✓ api/generator/sync_postman.py exists and is executable
✓ run_api_pipeline.sh has Step 15 (Postman generation)
✓ run_migration.ps1 exists
```

#### Step 1.2: Verify Prerequisites
```bash
# Check Python installation
python3 --version        # Should be 3.7+

# Check if pyyaml is installed
python3 -c "import yaml; print('✓ pyyaml installed')"

# Check npm/nodejs (for Newman)
npm --version           # Should be 6.0+
node --version          # Should be 12.0+

# Check curl (for Postman API calls)
curl --version          # Should be available
```

**If any are missing:**
```bash
# Install Python dependencies
pip install pyyaml requests

# Install Newman globally
npm install -g newman

# Or use npx (no installation needed)
npx newman run --version
```

#### Step 1.3: Create/Update .env File
```bash
cd E:\Projects\InsureTech

# Copy example to .env if not present
[ ! -f .env ] && cp .env.example .env

# Edit .env to add Postman settings
cat >> .env << 'EOF'

# ── Postman API Integration ───────────────────────────────────────────
# Optional: Set POSTMAN_API_KEY to auto-upload collections
# Get your key at: https://go.postman.co/settings/me/api-keys
POSTMAN_API_KEY=                    # Leave empty to skip upload
POSTMAN_WORKSPACE_ID=              # Your workspace ID (optional)
POSTMAN_COLLECTION_ID=             # Existing collection ID to update (optional)
EOF

# Verify .env is readable
grep POSTMAN .env
```

### Phase 2: Bash Pipeline Integration (10-15 minutes)

#### Step 2.1: Verify Existing Implementation
```bash
cd E:\Projects\InsureTech

# Check that Step 15 exists in run_api_pipeline.sh
grep -n "Step 15" run_api_pipeline.sh
grep -n "sync_postman.py" run_api_pipeline.sh

# Expected output:
# 296:step 15 16 "Generating Postman collection + environments..."
# 308:"$PY_CMD" sync_postman.py \
# 315:"$PY_CMD" sync_postman.py --upload \
```

**Status:** ✅ Already implemented at lines 296-331

#### Step 2.2: Verify Postman API Key Loading
```bash
# Check that POSTMAN_API_KEY is loaded from .env
grep -A5 "POSTMAN_API_KEY=" run_api_pipeline.sh

# Expected code pattern:
# if [ -f "$PROJECT_ROOT/.env" ]; then
#     POSTMAN_API_KEY=$(grep '^POSTMAN_API_KEY=' ...)
```

**Status:** ✅ Already implemented at lines 299-303

#### Step 2.3: Test Bash Script Execution
```bash
# Dry run without full pipeline
cd E:\Projects\InsureTech/api/generator

# Test sync_postman.py generation
python3 sync_postman.py

# Expected output:
# [1/4] Loading OpenAPI spec...
# [2/4] Building collection...
# [3/4] Generating environments...
# [4/4] Postman API sync...
# ✓ Postman Collection Ready!

# Verify generated files
ls -lh ../postman/InsureTech*.json
```

**Success Criteria:**
- ✅ `InsureTech.postman_collection.json` generated (>500KB)
- ✅ `InsureTech_local.postman_environment.json` generated
- ✅ `InsureTech_staging.postman_environment.json` generated
- ✅ `InsureTech_production.postman_environment.json` generated
- ✅ `InsureTech_newman_test.postman_environment.json` generated

### Phase 3: PowerShell Integration (15-20 minutes)

#### Step 3.1: Review Migration Script
```powershell
cd E:\Projects\InsureTech

# Check current structure of run_migration.ps1
Get-Content run_migration.ps1 | Select-String -Pattern "^# Step" | Head -10

# Expected output:
# # Step 0: Bootstrap prerequisites
# # Step 1: Find project root
# # Step 2: Load .env
# # etc.
```

#### Step 3.2: Add Postman Integration to run_migration.ps1

Two options:

**Option A: Include the patch (Recommended)**
```powershell
# At the end of run_migration.ps1, add:
& ".\POSTMAN_MIGRATION_PATCH.ps1" -ProjectRoot $PSScriptRoot

# This will run the Postman generation as a final step
```

**Option B: Inline the code**

Add this section after successful migration (around line 145):

```powershell
# Step 8: Generate Postman Collection
Write-Host "`n[8/8] Generating Postman collection..." -ForegroundColor Yellow

$generatorDir = Join-Path $projectRoot "api" "generator"
if (Test-Path $generatorDir) {
    Push-Location $generatorDir
    
    # Generate locally
    & python sync_postman.py
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "  ✓ Postman collection generated" -ForegroundColor Green
        
        # Upload if API key available
        if ($env:POSTMAN_API_KEY) {
            Write-Host "  ✓ Uploading to Postman API..." -ForegroundColor Cyan
            & python sync_postman.py --upload
        }
    } else {
        Write-Host "  ⚠ Postman generation skipped" -ForegroundColor Yellow
    }
    
    Pop-Location
}
```

#### Step 3.3: Test PowerShell Script
```powershell
cd E:\Projects\InsureTech

# Dry run the migration script
.\run_migration.ps1 -DryRun

# Expected output:
# [1/7] Checking prerequisites... ✓
# [2/7] Finding project root... ✓
# [3/7] Loading environment variables... ✓
# [4/7] Verifying configuration... ✓
# [5/7] Checking SSL certificates... ✓
# [6/7] Regenerating proto code... ✓
# [7/7] Running migration... (DRY RUN)
# [8/8] Generating Postman collection... (DRY RUN)
```

#### Step 3.4: Full Run (with real database)
```powershell
cd E:\Projects\InsureTech

# Full migration including Postman
.\run_migration.ps1

# Success indicators:
# - No errors during migration
# - Postman collection generated
# - If POSTMAN_API_KEY set: Collection uploaded to Postman API
```

### Phase 4: Enhanced Environment Propagation (10-15 minutes)

#### Step 4.1: Apply sync_postman.py Enhancements

The enhancements in `sync_postman_enhanced.py` are backward compatible. To apply:

**Option A: Manual Application**

1. Open `api/generator/sync_postman.py`
2. Follow the changes documented in `sync_postman_enhanced.py`
3. Focus on CHANGE 1: Enhanced make_test_script() with environment propagation

**Option B: Automated Script** (once all files are ready)

```bash
cd E:\Projects\InsureTech/api/generator
python3 apply_enhancements.py sync_postman.py
```

#### Step 4.2: Verify Environment Propagation

```bash
cd E:\Projects\InsureTech/api/generator

# Generate collection with enhanced scripts
python3 sync_postman.py

# Check that environment propagation is in test scripts
grep -c "pm.environment.set" ../postman/InsureTech.postman_collection.json

# Expected: Hundreds of environment.set calls (should be > 100)
```

### Phase 5: Testing and Validation (20-30 minutes)

#### Step 5.1: Import into Postman Desktop

```
1. Open Postman
2. Click "Import" (top left)
3. Select: api/postman/InsureTech.postman_collection.json
4. Click "Import"
5. Environment → Import → api/postman/InsureTech_local.postman_environment.json
6. Click environment dropdown → Select "InsureTech — Local"
7. In environment, set:
   - base_url = http://localhost:8080
   - user_mobile_number = +1-555-0100
   - login_device_type = ANDROID
```

#### Step 5.2: Test B2C OTP Flow in Postman

```
1. Open collection: AuthService → POST /v1/auth/otp:send
2. Click "Send"
3. Expected response: 200 OK, success=true, data.otp_id
4. Check environment — otp_id should be set automatically
5. Set mobile_otp_code manually (from SMS or test backend)
6. Run POST /v1/auth/otp:verify
7. Run POST /v1/auth/login
8. Expected: access_token and refresh_token captured automatically
9. Run any protected endpoint (e.g., GET /v1/users/me)
10. Should work without manual token setup!
```

**Success Criteria:**
- ✅ Token captured automatically from login response
- ✅ Bearer token injected in subsequent requests
- ✅ No "401 Unauthorized" errors
- ✅ All test scripts pass

#### Step 5.3: Test with Newman CLI

```bash
cd E:\Projects\InsureTech

# Run auth smoke tests
npx newman run api/postman/auth_smoke.postman_collection.json \
  -e api/postman/InsureTech_newman_test.postman_environment.json \
  --reporters cli,json \
  --reporter-json-export auth_smoke_results.json

# Expected output:
# → Running collection [InsureTech Auth Smoke]
# ✓ B2C Mobile OTP → JWT
#   ✓ 01 OTP Send (Mobile)
#   ✓ 02 OTP Verify (Mobile)
#   ✓ 03 Login After OTP
#   ... more tests ...
# ░ Collections: 1 of 1
# ✓ Requests: 20 of 20 (passing)
# ✓ Tests: 150 of 150 (passing)
```

#### Step 5.4: Test Authorization Enforcement

```bash
# Run authz tests
npx newman run api/postman/b2c_authz_enforcement.postman_collection.json \
  -e api/postman/InsureTech_local.postman_environment.json \
  --reporters cli

# Expected: All authorization tests pass
```

### Phase 6: CI/CD Integration (15-20 minutes)

#### Step 6.1: GitHub Actions Integration

Create `.github/workflows/postman-tests.yml`:

```yaml
name: Postman API Tests

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main, develop ]

jobs:
  postman-tests:
    runs-on: ubuntu-latest
    
    services:
      api:
        image: insuretech-api:latest
        ports:
          - 8080:8080
        env:
          DB_HOST: postgres
          DB_PORT: 5432
    
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '18'
      
      - name: Install Newman
        run: npm install -g newman
      
      - name: Wait for API
        run: |
          for i in {1..30}; do
            curl -sf http://localhost:8080/health && break
            sleep 1
          done
      
      - name: Run Postman Auth Smoke Tests
        run: |
          newman run api/postman/auth_smoke.postman_collection.json \
            -e api/postman/InsureTech_newman_test.postman_environment.json \
            --env-var base_url=http://localhost:8080 \
            --reporters cli,json,htmlextra \
            --reporter-json-export auth_results.json \
            --reporter-htmlextra-export auth_results.html
      
      - name: Run Full Collection Tests
        run: |
          newman run api/postman/InsureTech.postman_collection.json \
            -e api/postman/InsureTech_local.postman_environment.json \
            --env-var base_url=http://localhost:8080 \
            --reporters cli,json \
            --reporter-json-export full_results.json \
            --timeout-request 5000
      
      - name: Upload Test Results
        if: always()
        uses: actions/upload-artifact@v3
        with:
          name: postman-results
          path: |
            auth_results.json
            auth_results.html
            full_results.json
```

#### Step 6.2: GitLab CI Integration

Create `.gitlab-ci.yml` section:

```yaml
postman:tests:
  stage: test
  image: node:18
  services:
    - insuretech-api:latest
      alias: api
  script:
    - npm install -g newman
    - npm install -g newman-reporter-htmlextra
    - sleep 10  # Wait for API to start
    - newman run api/postman/auth_smoke.postman_collection.json
        -e api/postman/InsureTech_newman_test.postman_environment.json
        --env-var base_url=http://api:8080
        --reporters cli,json,htmlextra
        --reporter-json-export auth_results.json
        --reporter-htmlextra-export auth_results.html
  artifacts:
    when: always
    paths:
      - auth_results.json
      - auth_results.html
    reports:
      junit: auth_results.json
```

### Phase 7: Postman Cloud Integration (Optional, 10 minutes)

#### Step 7.1: Get Postman API Key

```
1. Sign in to https://go.postman.co
2. Go to: Settings → API Keys
3. Click "Generate API Key"
4. Copy the key
5. Add to .env file: POSTMAN_API_KEY=your-key-here
```

#### Step 7.2: Configure Workspace ID (Optional)

```
1. In Postman, click workspace selector
2. Click "Settings" next to your workspace
3. Note the Workspace ID from the URL or settings page
4. Add to .env: POSTMAN_WORKSPACE_ID=your-workspace-id
```

#### Step 7.3: Auto-Upload Collections

```bash
cd E:\Projects\InsureTech/api/generator

# This will upload to Postman cloud
POSTMAN_API_KEY=your-key python3 sync_postman.py --upload

# Expected output:
# ✓ Collection 'InsureTech API' synced to Postman
# ✓ Collection ID: abc123...
# ✓ View at: https://go.postman.co/collection/abc123...
# ✓ InsureTech — Local synced to Postman
# ✓ InsureTech — Staging synced to Postman
# ✓ InsureTech — Production synced to Postman
```

## Verification Checklist

### ✅ Bash Pipeline
- [ ] `run_api_pipeline.sh` has Step 15 for Postman
- [ ] POSTMAN_API_KEY is loaded from .env
- [ ] Collections are generated at `api/postman/`
- [ ] `--upload` flag works when API key is present
- [ ] All environment files are created

### ✅ PowerShell Pipeline
- [ ] `run_migration.ps1` includes Postman generation
- [ ] Script calls `sync_postman.py` after migration
- [ ] POSTMAN_API_KEY is passed to the generator
- [ ] Script handles missing API key gracefully
- [ ] Error handling is in place

### ✅ Environment Propagation
- [ ] Test scripts include `pm.environment.set()` calls
- [ ] Tokens are auto-captured from responses
- [ ] Resource IDs are captured for CRUD operations
- [ ] Both singular and `last_*` variants are set
- [ ] All auth flows propagate required variables

### ✅ Documentation
- [ ] `POSTMAN_COLLECTION_GUIDE.md` is complete
- [ ] All auth flows are documented
- [ ] Newman CLI examples provided
- [ ] CI/CD integration examples included
- [ ] Troubleshooting guide is available

### ✅ Testing
- [ ] Collections import without errors
- [ ] Environments are selectable in Postman
- [ ] Auth smoke tests pass locally
- [ ] Full collection tests pass
- [ ] Newman tests work from CI/CD

### ✅ Cloud Integration (Optional)
- [ ] POSTMAN_API_KEY is set in .env
- [ ] Collections upload to Postman successfully
- [ ] Collections visible at https://go.postman.co/collections
- [ ] Workspace ID is correct (if specified)
- [ ] Collection updates work for existing collections

## Troubleshooting Common Issues

### Issue: "ModuleNotFoundError: No module named 'yaml'"
```bash
# Solution: Install pyyaml
pip install pyyaml
```

### Issue: "Newman not found"
```bash
# Solution: Install Newman globally
npm install -g newman

# Or use npx (no install)
npx newman run --version
```

### Issue: "POSTMAN_API_KEY rejected by server"
```bash
# Check:
1. Key is valid (test at https://api.getpostman.com)
2. Key has proper permissions
3. Workspace ID is correct (if specified)
4. Collection ID exists (if updating existing)

# Test with curl:
curl -X GET https://api.getpostman.com/me \
  -H "X-API-Key: your-api-key"
```

### Issue: "Collections not generating"
```bash
# Check:
1. OpenAPI spec exists: api/openapi.yaml
2. Generator dependencies installed: pip install -r api/generator/requirements.txt
3. Python version >= 3.7
4. Run with verbose output: python -u sync_postman.py
```

### Issue: "Tests failing in Newman but passing in Postman"
```bash
# Common causes:
1. Base URL different (check --env-var base_url=...)
2. Timeouts different (use --timeout-request 10000)
3. Authentication not working (check environment variables)
4. OTP codes expired (use fresh codes)

# Debug:
newman run ... -v  # Verbose output
newman run ... --reporters cli,json --reporter-json-export debug.json
```

## Next Steps

1. **Immediate** (Today)
   - [ ] Review all documentation
   - [ ] Run sync_postman.py locally
   - [ ] Import collection into Postman
   - [ ] Test B2C OTP flow

2. **Short-term** (This week)
   - [ ] Add to run_migration.ps1
   - [ ] Test PowerShell integration
   - [ ] Set up Postman API key (optional)
   - [ ] Run Newman tests locally

3. **Medium-term** (This sprint)
   - [ ] Add to GitHub Actions
   - [ ] Configure GitLab CI
   - [ ] Upload to Postman cloud
   - [ ] Document for team

4. **Long-term** (Ongoing)
   - [ ] Monitor test pass rates
   - [ ] Update collections on API changes
   - [ ] Expand auth test coverage
   - [ ] Add performance testing

## Support and Resources

- **Postman Learning Center:** https://learning.postman.com/
- **Newman Documentation:** https://github.com/postmanlabs/newman
- **InsureTech API Documentation:** `docs/` or `https://localhost:8080/docs`
- **Team Support:** Ask @dev-team in #api-testing channel

---

**Document Version:** 1.0  
**Last Updated:** 2024-01-15  
**Maintained By:** API Team  

