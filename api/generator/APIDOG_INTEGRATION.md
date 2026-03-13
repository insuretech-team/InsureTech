# Apidog Integration

This document explains how the Apidog integration works for automatic API testing and mocking.

## Overview

The pipeline automatically syncs the OpenAPI specification, schemas, endpoints, DTOs, enums, and paths to Apidog for:
- **API Mocking** - Auto-generated mock servers for all endpoints
- **API Testing** - Automated test generation
- **API Documentation** - Interactive API testing interface
- **Collaboration** - Team access to live API specs

## Setup

### 1. Get Your Apidog Token

1. Log in to [Apidog](https://apidog.com)
2. Go to **Settings** → **API Keys**
3. Generate a new API token
4. Copy the token (format: `APS-xxxxxxxxxxxxx`)

### 2. Configure Local Environment

Add the token to your `.env` file in the project root:

```env
API_DOG_TOKEN=APS-xxxxxxxxxxxxx
```

**Note:** The `.env` file is already in `.gitignore` - never commit your token!

### 3. Configure GitHub Actions (Optional)

For CI/CD integration:

1. Go to your GitHub repository
2. Navigate to **Settings** → **Secrets and variables** → **Actions**
3. Click **New repository secret**
4. Name: `API_DOG_TOKEN`
5. Value: Your Apidog token
6. Click **Add secret**

## Usage

### Local Pipeline

The Apidog sync is automatically triggered when you run:

```powershell
.\run_api_pipeline.ps1
```

**What gets synced:**
- ✅ OpenAPI 3.1 specification (`openapi.yaml`)
- ✅ All schemas (`schemas/**/*.yaml`)
- ✅ All enums (`enums/**/*.yaml`)
- ✅ All endpoints (`paths/**/*.yaml`)
- ✅ Mock server configuration

**Output:**
```
[13/15] Syncing to Apidog...
  Found API_DOG_TOKEN, syncing to Apidog...
  
  ============================================================
    Apidog API Sync
  ============================================================
  
  ✓ Connected to Apidog - Project: InsureTech Platform
  
  [1/5] Uploading OpenAPI Specification...
  ✓ OpenAPI spec uploaded successfully
    - Endpoints: 221
    - Schemas: 740
  
  [2/5] Syncing Schemas...
  ✓ Synced 740 schemas
  
  [3/5] Syncing Enums...
  ✓ Synced 125 enums
  
  [4/5] Syncing Endpoints...
  ✓ Synced 221 endpoints
  
  [5/5] Enabling Mock Servers...
  ✓ Mock server enabled: https://mock.apidog.com/xxxxx
  
  ============================================================
    Sync Summary
  ============================================================
  OpenAPI Spec: ✓ Uploaded
  Schemas:      740 synced
  Enums:        125 synced
  Endpoints:    221 synced
  Mock Servers: ✓ Enabled
  ============================================================
```

### CI/CD Pipeline

The sync automatically runs on GitHub Actions when:
- Pushing to `main` branch
- OpenAPI spec or proto files change

**Workflow step:**
```yaml
- name: Sync to Apidog
  if: github.ref == 'refs/heads/main' && github.event_name == 'push'
  env:
    API_DOG_TOKEN: ${{ secrets.API_DOG_TOKEN }}
  run: |
    if [ -n "$API_DOG_TOKEN" ]; then
      echo "Syncing to Apidog..."
      cd api/generator
      pip install requests pyyaml
      python sync_apidog.py || echo "Apidog sync failed, continuing..."
    else
      echo "API_DOG_TOKEN not found, skipping Apidog sync"
    fi
  continue-on-error: true
```

### Manual Sync

You can also run the sync script manually:

```powershell
cd api/generator
python sync_apidog.py
```

Or with explicit token:

```powershell
$env:API_DOG_TOKEN="APS-xxxxxxxxxxxxx"
python sync_apidog.py
```

## Features

### 1. Mock Server

Apidog automatically creates a mock server for all endpoints:
- Auto-generated responses based on OpenAPI schema
- Realistic data generation
- Configurable response delay (default: 100ms)

**Access:**
- Mock URL: `https://mock.apidog.com/your-project-id`
- Use in development/testing environments

### 2. Automated Tests

Apidog generates tests for:
- Schema validation
- Response status codes
- Required fields
- Data type validation
- Enum value validation

### 3. Interactive Documentation

- Try out endpoints directly from browser
- View request/response examples
- Test authentication flows
- Share with team members

## Script Details

### `sync_apidog.py`

**Location:** `api/generator/sync_apidog.py`

**Dependencies:**
```python
requests  # HTTP client
pyyaml    # YAML parsing
```

**Main Functions:**

```python
class ApidogSync:
    def test_connection() -> bool
        # Test API connection and get project ID
    
    def upload_openapi_spec(spec_path: Path) -> bool
        # Upload OpenAPI 3.1 spec
    
    def sync_schemas(schemas_dir: Path) -> int
        # Sync all schema definitions
    
    def sync_enums(enums_dir: Path) -> int
        # Sync all enum definitions
    
    def sync_paths(paths_dir: Path) -> int
        # Sync all endpoint definitions
    
    def create_mock_servers() -> bool
        # Enable mock servers
```

**Exit Codes:**
- `0` - Success
- `1` - Failed to upload OpenAPI spec (critical)

## Troubleshooting

### Token Not Found

**Error:**
```
✗ Error: API_DOG_TOKEN environment variable not set
```

**Solution:**
1. Check `.env` file exists in project root
2. Verify token format: `API_DOG_TOKEN=APS-xxxxxxxxxxxxx`
3. Restart PowerShell session

### Connection Failed

**Error:**
```
✗ Failed to connect to Apidog: 401
```

**Solution:**
1. Verify token is valid and not expired
2. Check internet connection
3. Ensure token has correct permissions

### Upload Failed

**Error:**
```
✗ Upload failed: 422 - Invalid OpenAPI spec
```

**Solution:**
1. Run OpenAPI validation: `python enhanced_validator.py ../openapi.yaml`
2. Fix any critical errors
3. Regenerate spec: `.\run_api_pipeline.ps1`

### Sync Skipped

**Output:**
```
⚠ API_DOG_TOKEN not found, skipping Apidog sync
```

**This is normal if:**
- You haven't set up Apidog yet
- Running in CI without the secret configured
- Token is optional for local development

## API Documentation

### Apidog REST API

**Base URL:** `https://api.apidog.com/api/v1`

**Authentication:**
```http
Authorization: Bearer APS-xxxxxxxxxxxxx
```

**Endpoints Used:**

| Endpoint | Method | Purpose |
|----------|--------|---------|
| `/projects` | GET | List projects, get project ID |
| `/imports/openapi` | POST | Import OpenAPI specification |
| `/schemas` | POST | Create/update schema definitions |
| `/types` | POST | Create/update enum types |
| `/endpoints` | POST | Create/update API endpoints |
| `/mock-servers` | POST | Enable mock servers |

**Rate Limits:**
- 100 requests per minute (per token)
- Bulk upload operations count as 1 request

## Best Practices

### 1. Separate Tokens for Environments

```env
# Development
API_DOG_TOKEN=APS-dev-xxxxx

# Staging
API_DOG_TOKEN=APS-staging-xxxxx

# Production (read-only)
API_DOG_TOKEN=APS-prod-xxxxx
```

### 2. Use Mock Servers in Tests

```javascript
// Use Apidog mock instead of real backend
const API_BASE = process.env.CI 
  ? 'https://mock.apidog.com/xxxxx'
  : 'http://localhost:8080';
```

### 3. Keep Tokens Secure

- ✅ Store in `.env` file (gitignored)
- ✅ Use GitHub Secrets for CI/CD
- ✅ Rotate tokens every 90 days
- ❌ Never commit tokens to repository
- ❌ Never share tokens in Slack/email

### 4. Monitor Sync Status

Check GitHub Actions logs for:
- Successful syncs
- Failed uploads
- API changes detected

## Related Files

- `sync_apidog.py` - Main sync script
- `run_api_pipeline.ps1` - Pipeline with Apidog integration
- `.github/workflows/openapi-validation.yml` - CI/CD workflow
- `.env` - Local token storage (gitignored)
- `openapi.yaml` - Source OpenAPI spec

## Support

- **Apidog Documentation:** https://apidog.com/docs
- **API Reference:** https://api.apidog.com/docs
- **Support:** support@apidog.com

---

**Last Updated:** January 2026  
**Integration Version:** 1.0
