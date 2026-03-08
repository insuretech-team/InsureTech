# Nginx Scripts Directory

This directory contains automation scripts for managing the nginx configuration.

---

## Scripts Overview

### 1. `validate-config.ps1` ✅
**Purpose:** Validate nginx configuration files  
**Platform:** Windows (PowerShell)  
**Usage:**
```powershell
.\validate-config.ps1
```
**What it does:**
- Checks all 42 configuration requirements
- Validates file presence
- Checks content patterns
- Reports validation results

**When to use:** Before any deployment, after configuration changes

---

### 2. `enable-sites.ps1` / `enable-sites.sh` ✅
**Purpose:** Enable nginx virtual host sites  
**Platform:** Windows (PowerShell) / Linux (Bash)  
**Usage:**
```powershell
# Windows
.\enable-sites.ps1

# Linux
./enable-sites.sh
```
**What it does:**
- Creates symlinks/copies from sites-available to sites-enabled
- Enables all 4 configured sites
- Verifies successful enablement

**When to use:** Initial setup, when adding new sites

---

### 3. `create-dev-branch.ps1` ✅
**Purpose:** Create development branch for nginx migration  
**Platform:** Windows (PowerShell)  
**Usage:**
```powershell
.\create-dev-branch.ps1
```
**What it does:**
- Creates `dev/nginx-modular` branch
- Stages all nginx-related changes
- Commits with detailed message
- Pushes to remote (optional)

**When to use:** Once, before starting dev deployment

---

### 4. `pre-deployment-check.ps1` ✅
**Purpose:** Comprehensive pre-deployment validation  
**Platform:** Windows (PowerShell)  
**Usage:**
```powershell
.\pre-deployment-check.ps1
```
**What it does:**
- Validates all configuration files
- Checks Docker setup
- Verifies CI/CD workflows
- Tests backward compatibility
- Optional: Docker build test

**When to use:** Before creating dev branch, before production deployment

---

### 5. `setup.sh` 📋
**Purpose:** Production server setup  
**Platform:** Linux (Bash)  
**Usage:**
```bash
sudo ./setup.sh
```
**What it does:**
- Copies nginx configuration to /etc/nginx/
- Creates required directories
- Sets proper permissions
- Enables sites
- Validates configuration

**When to use:** Initial production setup (if not using Docker)

---

### 6. `test-config.sh` 📋
**Purpose:** Test nginx configuration syntax  
**Platform:** Linux (Bash)  
**Usage:**
```bash
./test-config.sh
```
**What it does:**
- Runs `nginx -t` to validate syntax
- Shows validation results
- Reports any errors

**When to use:** After configuration changes, before reload

---

### 7. `clear-cache.sh` 📋
**Purpose:** Clear nginx cache  
**Platform:** Linux (Bash)  
**Usage:**
```bash
sudo ./clear-cache.sh
```
**What it does:**
- Removes all cached files
- Clears static, API, and microcache
- Reloads nginx

**When to use:** When cache needs to be purged, troubleshooting

---

## Script Execution Order

### First Time Setup
1. `validate-config.ps1` - Validate configuration
2. `enable-sites.ps1` - Enable sites
3. `pre-deployment-check.ps1` - Pre-deployment validation
4. `create-dev-branch.ps1` - Create dev branch

### Regular Operations
- `validate-config.ps1` - After any config changes
- `test-config.sh` - Before nginx reload
- `clear-cache.sh` - When cache needs clearing

---

## Platform Notes

### Windows (PowerShell)
- Run scripts from Trendico root directory
- May require execution policy adjustment:
  ```powershell
  Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
  ```

### Linux (Bash)
- Make scripts executable:
  ```bash
  chmod +x *.sh
  ```
- Some scripts require sudo

---

## Troubleshooting

### Script Won't Run (Windows)
```powershell
# Check execution policy
Get-ExecutionPolicy

# Set to RemoteSigned
Set-ExecutionPolicy RemoteSigned -Scope CurrentUser
```

### Script Won't Run (Linux)
```bash
# Make executable
chmod +x script-name.sh

# Check shebang
head -1 script-name.sh  # Should be #!/bin/bash
```

### "Command not found" (Linux)
```bash
# Run with explicit path
./script-name.sh

# Or add to PATH
export PATH=$PATH:$(pwd)
```

---

## Quick Reference

| Task | Command |
|------|---------|
| Validate config | `.\validate-config.ps1` |
| Enable sites | `.\enable-sites.ps1` |
| Pre-deployment check | `.\pre-deployment-check.ps1` |
| Create dev branch | `.\create-dev-branch.ps1` |
| Test syntax (Linux) | `./test-config.sh` |
| Clear cache (Linux) | `sudo ./clear-cache.sh` |

---

## Related Documentation

- **MIGRATION_GUIDE.md** - Complete migration steps
- **DEPLOYMENT_CHECKLIST.md** - Deployment procedures
- **QUICKSTART.md** - Quick start guide
- **README.md** - Main documentation

---

**Last Updated:** January 24, 2025
