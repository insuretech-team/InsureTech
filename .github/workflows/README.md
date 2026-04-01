# GitHub Workflows Documentation

## Overview

This directory contains GitHub Actions workflows for automated SDK publishing, validation, and version management.

## Workflows

### 1. `auto-publish-sdks.yml` - Automatic SDK Publishing

**Purpose**: Automatically publishes both Go and TypeScript SDKs when changes are detected.

**Triggers**:
- **Push to main**: When SDK files, API specs, or proto files change
- **Pull Request**: When PR is merged to main
- **Manual Dispatch**: Manual trigger with custom options

**Features**:
- ✅ First-time publication detection (0.1.0 for new packages)
- ✅ Automatic version bumping based on commit messages
- ✅ Semantic versioning (major.minor.patch)
- ✅ Detects if repo/package exists before publishing
- ✅ Creates GitHub releases automatically
- ✅ Publishes to NPM (TypeScript) and GitHub (Go)

**Version Bump Detection**:
- **MAJOR**: Commit message contains "BREAKING CHANGE" or "major:"
- **MINOR**: Commit message starts with "feat:", "feature:", or "minor:"
- **PATCH**: All other commits (default)

**Manual Trigger Options**:
```yaml
version_bump: [patch, minor, major]
sdk_to_publish: [both, go, typescript]
skip_tests: [true, false]
```

**Usage Examples**:

```bash
# Automatic trigger - commit with version hint
git commit -m "feat: Add new payment service" # Minor bump
git commit -m "fix: Correct policy validation" # Patch bump
git commit -m "major: Remove deprecated endpoints

BREAKING CHANGE: Removed old API v1" # Major bump

# Manual trigger via GitHub UI
# Go to Actions > Auto-Publish SDKs > Run workflow
# Select options and run
```

### 2. `publish-go-sdk.yml` - Go SDK Publication

**Purpose**: Publishes Go SDK to separate repository (newage-saint/insuretech-go-sdk)

**Triggers**:
- Tag push: `go-v*.*.*`
- Manual dispatch with version input

**Target**: https://github.com/newage-saint/insuretech-go-sdk

### 3. `publish-typescript-sdk.yml` - TypeScript SDK Publication

**Purpose**: Publishes TypeScript SDK to NPM

**Triggers**:
- Tag push: `ts-v*.*.*`
- Manual dispatch with version input

**Target**: https://www.npmjs.com/package/@newage-saint/insuretech-sdk

### 4. `openapi-validation.yml` - API Validation

**Purpose**: Validates OpenAPI specifications and generates documentation

**Triggers**:
- Push to main
- Pull requests to main
- Changes to API files

## Version Management

### Semantic Versioning

We follow [Semantic Versioning 2.0.0](https://semver.org/):

- **MAJOR** (X.0.0): Breaking changes
- **MINOR** (0.X.0): New features (backwards compatible)
- **PATCH** (0.0.X): Bug fixes (backwards compatible)

### First-Time Publication

When SDKs are published for the first time:
- **Initial Version**: `0.1.0`
- Workflow detects if repository/package exists
- Creates repository if needed (Go SDK)
- Publishes to NPM with first-time messaging (TypeScript SDK)

### Versioning Strategy

#### Go SDK
- Versions stored as Git tags: `v1.0.0`, `v1.1.0`, etc.
- Tags fetched from target repository (newage-saint/insuretech-go-sdk)
- If no tags exist: starts at `0.1.0`

#### TypeScript SDK
- Version stored in `package.json` and NPM
- Checked via `npm view @newage-saint/insuretech-sdk version`
- If package doesn't exist: starts at `0.1.0`

## Required Secrets

Configure these secrets in GitHub repository settings:

| Secret | Description | Used By |
|--------|-------------|---------|
| `NEWAGE_SAINT_PAT` | Personal Access Token for newage-saint account | Go SDK publishing |
| `NPM_TOKEN` | NPM authentication token | TypeScript SDK publishing |
| `GITHUB_TOKEN` | Automatically provided by GitHub | Release creation |

### Setting Up Secrets

1. **NEWAGE_SAINT_PAT**:
   ```bash
   # Create a Personal Access Token at:
   # https://github.com/settings/tokens
   # Required scopes: repo, workflow
   ```

2. **NPM_TOKEN**:
   ```bash
   # Create at: https://www.npmjs.com/settings/[username]/tokens
   # Type: Automation token
   npm token create --type automation
   ```

## Workflow Paths

Files that trigger automatic publishing:

```yaml
- 'sdks/insuretech-go-sdk/**'        # Go SDK changes
- 'sdks/insuretech-typescript-sdk/**' # TypeScript SDK changes
- 'api/openapi.yaml'                  # API specification
- 'proto/**'                          # Protocol buffer definitions
- 'buf.yaml'                          # Buf configuration
- 'buf.gen.yaml'                      # Buf generation config
```

## Monitoring & Debugging

### View Workflow Runs
1. Go to **Actions** tab in GitHub
2. Select the workflow
3. View run history and logs

### Common Issues

#### 1. Repository Doesn't Exist
**Error**: `Repository newage-saint/insuretech-go-sdk not found`

**Solution**: Workflow will automatically create it on first run

#### 2. NPM Package Doesn't Exist
**Error**: Package not found on NPM

**Solution**: Workflow will publish as first-time (v0.1.0)

#### 3. Permission Denied
**Error**: `Permission denied (publickey)`

**Solution**: Check that `NEWAGE_SAINT_PAT` secret is correctly configured

#### 4. Version Already Exists
**Error**: Tag already exists or NPM version conflict

**Solution**: Version bumping is automatic - ensure commits are following convention

## Best Practices

### Commit Messages

Use conventional commits for automatic version detection:

```bash
# Features (minor bump)
feat: Add claim submission endpoint
feature: Implement payment gateway

# Fixes (patch bump)
fix: Correct date parsing in policy service
patch: Update validation logic

# Breaking Changes (major bump)
major: Remove deprecated v1 API

BREAKING CHANGE: This removes all v1 endpoints
```

### Testing Before Release

1. **Create a feature branch**
2. **Make changes and commit**
3. **Create PR to main**
4. **Review CI checks**
5. **Merge to main** (triggers auto-publish)

### Manual Publishing

When you need control:

1. Go to **Actions** > **Auto-Publish SDKs on Main Commit**
2. Click **Run workflow**
3. Select:
   - Version bump type (patch/minor/major)
   - SDK to publish (both/go/typescript)
   - Skip tests (if needed, use with caution)
4. Click **Run workflow**

## Release Process

### Automatic Flow

```
Code Change → Push to Main → Detect Changes → Version Bump → 
Build & Test → Publish → Create Release → Tag
```

### First-Time Flow

```
Code Change → Push to Main → Detect No Repo/Package → Create v0.1.0 →
Build & Test → Create Repo (if needed) → Publish → Create Release
```

## Troubleshooting

### Enable Debug Logging

Add these secrets to enable verbose logging:
- `ACTIONS_STEP_DEBUG`: `true`
- `ACTIONS_RUNNER_DEBUG`: `true`

### Validate Workflow Locally

```bash
# Install act (GitHub Actions local runner)
# https://github.com/nektos/act

# Test workflow locally
act push -W .github/workflows/auto-publish-sdks.yml
```

### Check Workflow Syntax

```bash
# Install actionlint
# https://github.com/rhysd/actionlint

actionlint .github/workflows/auto-publish-sdks.yml
```

## Version History

| Version | Date | Changes |
|---------|------|---------|
| 2.0 | 2026-01-03 | Added first-time publication detection, manual triggers |
| 1.0 | 2025-12-XX | Initial automated workflow setup |

## Support

For issues or questions:
- Create an issue in the repository
- Check workflow run logs in Actions tab
- Review this documentation

---

**Last Updated**: January 3, 2026
**Maintained By**: InsureTech Platform Team
