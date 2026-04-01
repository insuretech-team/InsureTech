# Automated Versioning Guide

## How It Works

When you commit to `main`, the CI/CD automatically:
1. Detects what changed (Go SDK, TypeScript SDK, or API)
2. Determines version bump type from commit message
3. Calculates new version
4. Publishes SDKs with new version
5. Creates GitHub releases

---

## Version Bump Rules

The version bump type is determined by your commit message:

### MAJOR Version (1.0.0 → 2.0.0)
Breaking changes that require user code updates.

**Commit message keywords:**
- `BREAKING CHANGE:` (anywhere in commit)
- `major:` (at start of commit message)

**Example:**
```bash
git commit -m "major: Redesign authentication API

BREAKING CHANGE: Auth methods now require new parameters"
```

### MINOR Version (1.0.0 → 1.1.0)
New features that are backward compatible.

**Commit message keywords:**
- `feat:` (at start)
- `feature:` (at start)
- `minor:` (at start)

**Example:**
```bash
git commit -m "feat: Add new fraud detection endpoints"
```

### PATCH Version (1.0.0 → 1.0.1)
Bug fixes and minor changes (default).

**Any other commit message:**

**Example:**
```bash
git commit -m "Fix policy creation bug"
git commit -m "Update documentation"
git commit -m "Refactor service code"
```

---

## Commit Message Format (Recommended)

Follow Conventional Commits format:

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Types:
- `feat:` - New feature (MINOR bump)
- `fix:` - Bug fix (PATCH bump)
- `docs:` - Documentation only (PATCH bump)
- `style:` - Formatting (PATCH bump)
- `refactor:` - Code refactoring (PATCH bump)
- `test:` - Adding tests (PATCH bump)
- `chore:` - Maintenance (PATCH bump)
- `major:` - Breaking change (MAJOR bump)

### Examples:

```bash
# MAJOR bump
git commit -m "major: Redesign policy API

BREAKING CHANGE: PolicyCreationRequest now requires tenantId"

# MINOR bump
git commit -m "feat: Add voice commands service"
git commit -m "feature(analytics): Add real-time dashboards"

# PATCH bump (default)
git commit -m "fix: Correct payment verification logic"
git commit -m "docs: Update README examples"
git commit -m "refactor: Improve error handling"
```

---

## What Triggers Publication

### Go SDK is published when:
- Changes in `sdks/insuretech-go-sdk/`
- Changes in `api/openapi.yaml`
- Changes in `proto/`

### TypeScript SDK is published when:
- Changes in `sdks/insuretech-typescript-sdk/`
- Changes in `api/openapi.yaml`
- Changes in `proto/`

### Both SDKs are published when:
- Changes in `api/openapi.yaml` or `proto/` (API changes affect both)

---

## Version Management

### Starting Versions

If no tags exist:
- Initial version: `0.1.0`
- First stable: `1.0.0` (when you're ready)

### Version History

Versions are tracked via git tags:
- Go SDK: `go-v1.0.0`, `go-v1.0.1`, `go-v1.1.0`, etc.
- TypeScript SDK: `ts-v1.0.0`, `ts-v1.0.1`, `ts-v1.1.0`, etc.

### Check Current Versions

```bash
# Go SDK
git tag -l "go-v*" | sort -V | tail -n1

# TypeScript SDK  
git tag -l "ts-v*" | sort -V | tail -n1
```

---

## Workflow Behavior

### Scenario 1: API Changes
```bash
# Update OpenAPI spec
git add api/openapi.yaml
git commit -m "feat: Add new analytics endpoints"
git push origin main
```

**Result:**
- ✅ Both Go and TypeScript SDKs published
- Version bump: MINOR (1.0.0 → 1.1.0)
- Both get same version number

### Scenario 2: Go SDK Only
```bash
# Update Go SDK
git add sdks/insuretech-go-sdk/
git commit -m "fix: Improve retry logic in Go SDK"
git push origin main
```

**Result:**
- ✅ Only Go SDK published
- Version bump: PATCH (1.0.0 → 1.0.1)
- TypeScript SDK unchanged

### Scenario 3: TypeScript SDK Only
```bash
# Update TypeScript SDK
git add sdks/insuretech-typescript-sdk/
git commit -m "feat: Add request interceptors"
git push origin main
```

**Result:**
- ✅ Only TypeScript SDK published
- Version bump: MINOR (1.0.0 → 1.1.0)
- Go SDK unchanged

### Scenario 4: Breaking Change
```bash
git commit -m "major: Redesign authentication

BREAKING CHANGE: All auth methods now require tenantId parameter"
git push origin main
```

**Result:**
- ✅ Both SDKs published (if API changed)
- Version bump: MAJOR (1.5.3 → 2.0.0)

---

## Manual Override

If you need specific version numbers, you can still use the release scripts:

```powershell
# Force specific version for Go SDK
.\scripts\release-go-sdk.ps1 -Version "2.5.0"

# Force specific version for TypeScript SDK
.\scripts\release-typescript-sdk.ps1 -Version "2.5.0"
```

---

## Monitoring

### View Workflow Status
```
https://github.com/[your-org]/InsureTech/actions
```

### Check Published Versions

**Go SDK:**
```bash
go list -m -versions github.com/newage-saint/insuretech-go-sdk
```

**TypeScript SDK:**
```bash
npm view @newage-saint/insuretech-sdk versions
```

---

## Best Practices

1. ✅ **Use conventional commits** - Clear version bumps
2. ✅ **Commit related changes together** - Atomic releases
3. ✅ **Test locally first** - Before pushing to main
4. ✅ **Document breaking changes** - In commit body
5. ✅ **Review before merge** - Use PRs for main

---

**Updated:** January 3, 2026
