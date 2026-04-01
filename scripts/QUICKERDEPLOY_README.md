# quickerdeploy.sh - Faster Deployment Script

## Overview
`quickerdeploy.sh` is an optimized version of `quick_deploy.sh` that dramatically speeds up deployments by:
1. **Building Go binaries natively in WSL** (not inside Docker containers)
2. **Using slim alpine Dockerfiles** with pre-built binaries
3. **Supporting selective service builds** via `--services` flag
4. **Running builds in parallel** for maximum speed

## Usage

### Full deployment (all services)
```bash
bash scripts/quickerdeploy.sh
```

### Selective deployment (only specific services)
```bash
bash scripts/quickerdeploy.sh --services=authn,gateway,b2b_portal
```

### Skip build, just restart containers
```bash
bash scripts/quickerdeploy.sh --no-build
```

### Update nginx + SSL certs only
```bash
bash scripts/quickerdeploy.sh --nginx-only
```

## What Each Flag Does

### `--services=authn,gateway,b2b_portal` (comma-separated, no spaces)
- Only builds the specified services
- Supported services: `authn`, `authz`, `b2b`, `gateway`, `storage`, `dbops`, `b2b_portal`
- If omitted, all services are built

### `--no-build`
- **Skips STEP 1-3** (no binary compilation, no Docker build)
- Skips image tarball creation and transfer
- Only transfers docker-compose file and public assets
- Immediately restarts containers with already-loaded images on server
- **Use when**: Only environment variables changed, or testing container restarts

### `--nginx-only`
- **Skips all Docker operations**
- Only syncs nginx configs and updates SSL certificates
- Uses `force-recreate` to apply any .env changes to gateway/b2b_portal
- **Use when**: Only nginx/SSL changes needed

## Speed Improvements

The key speed improvement comes from building Go binaries **outside Docker**:

### Traditional approach (quick_deploy.sh):
```
Docker container starts → Downloads Go modules → Compiles source → Outputs binary → Creates alpine image
```

### New approach (quickerdeploy.sh):
```
WSL: Downloads modules & compiles (parallel for all services)
↓
Docker: Copies pre-built binary → Creates alpine image
```

This avoids the overhead of:
- Starting multiple large golang containers
- Re-downloading Go modules for each service
- Compiling inside slow Docker layer caching

## Build Steps

### STEP 0: Sync .env.prod
- Always happens first
- Ensures db-migrate gets correct Neon credentials

### STEP 1: Clean binaries
- Removes `bin/` directory
- Forces fresh native compilation

### STEP 2: Build Go binaries natively (PARALLEL)
- Runs for each selected Go service:
  - `authn` → `./backend/inscore/microservices/authn/cmd/server/main.go` → `bin/authn/server`
  - `authz` → `./backend/inscore/microservices/authz/cmd/server/main.go` → `bin/authz/server`
  - `b2b` → `./backend/inscore/microservices/b2b/cmd/server/main.go` → `bin/b2b/server`
  - `gateway` → `./backend/inscore/cmd/gateway/main.go` → `bin/gateway/server`
  - `storage` → `./backend/inscore/cmd/storage/main.go` → `bin/storage/server`
  - `dbops` → `./backend/inscore/cmd/dbops/main.go` → `bin/dbops/dbops`

- Build flags: `CGO_ENABLED=0 GOOS=linux GOARCH=amd64 -ldflags="-w -s"`
- All builds run in parallel using `&` and `wait`

### STEP 3: Build Docker images
- **Go services**: Use slim inline Dockerfiles (alpine:latest + pre-built binary)
- **b2b_portal**: Uses original multi-stage Node.js build (npm install + next build)
- **Infra services**: Pulls redis:7-alpine and apache/kafka:latest

### STEP 4: Save images to tarball
- Only saves selected service images (not all)
- Dramatically reduces transfer time for selective deployments

### STEP 5: Transfer and deploy
- Transfers docker-compose.yml and image tarball to remote
- Transfers b2b_portal public assets
- Restarts Docker stack, updates nginx, and manages SSL certs

## Slim Dockerfile Example

For Go services, the script uses inline heredoc Dockerfiles:

```dockerfile
FROM alpine:latest
WORKDIR /app
COPY bin/authn/server .
COPY .env .
COPY go.mod .
COPY ops ./ops
COPY backend/inscore/configs ./backend/inscore/configs
COPY backend/inscore/templates ./backend/inscore/templates
COPY backend/inscore/secrets ./backend/inscore/secrets
CMD ["./server"]
```

These are **much smaller and faster** than golang:1.25-alpine builder images because they:
- Don't include the Go compiler
- Don't need to download/compile anything
- Start from minimal alpine base

## Expected Timing

On a modern WSL system with good bandwidth:

| Step | Time | Notes |
|------|------|-------|
| STEP 1: Clean | ~1s | Fast |
| STEP 2: Go builds (parallel) | ~30-60s | Depends on service count & CPU |
| STEP 3: Docker builds | ~20-30s | Go services are tiny; b2b_portal takes longer |
| STEP 4: Save to tarball | ~5-10s | Small images = small tarball |
| STEP 5: Transfer + Deploy | ~2-5 min | Network-dependent |
| **Total** | **~3-6 min** | vs ~15-20 min for quick_deploy.sh |

## Parallel Build Times

When building multiple services in STEP 2:
- **Single service**: ~10-15s
- **All 6 Go services**: ~30-60s (runs in parallel)

This is why `--services=authn` is much faster than full build.

## Common Workflows

### Quick fix to authn service
```bash
bash scripts/quickerdeploy.sh --services=authn
# Takes ~1-2 minutes total
```

### Update multiple services
```bash
bash scripts/quickerdeploy.sh --services=authn,gateway,authz
# Takes ~2-3 minutes total
```

### Full deployment with all services
```bash
bash scripts/quickerdeploy.sh
# Takes ~5-6 minutes total
```

### Restart after env change
```bash
bash scripts/quickerdeploy.sh --no-build
# Takes ~1-2 minutes (no compilation)
```

### Just update nginx configs
```bash
bash scripts/quickerdeploy.sh --nginx-only
# Takes ~30-60 seconds
```

## Requirements

- **WSL2 with Ubuntu** and Go 1.25+ installed
- **Docker** (via Docker Desktop WSL integration)
- **SSH access** to remote server (insureadmin@146.190.97.242)
- **.env.prod** file in project root with `B2B_PORTAL_API_BASE_URL`

## Troubleshooting

### Build fails with "go build" error
- Ensure Go 1.25+ is installed in WSL: `go version`
- Check that paths in SERVICE_PATHS match your actual directory structure

### Docker image not found on remote
- Verify tarball was transferred: `ssh insureadmin@146.190.97.242 ls -lh /home/insureadmin/insuretech/insuretech_images.tar`
- Check docker load worked: `ssh insureadmin@146.190.97.242 docker images | grep insuretech`

### Container not starting
- Check logs: `docker compose --profile full logs <service>`
- Verify .env.prod was synced: `ssh insureadmin@146.190.97.242 cat /home/insureadmin/insuretech/.env.prod`

## Comparison with quick_deploy.sh

| Feature | quick_deploy.sh | quickerdeploy.sh |
|---------|-----------------|-----------------|
| Build location | Inside Docker | Native WSL |
| Parallel builds | No | Yes |
| Service selection | No | Yes (--services) |
| Slim images | No | Yes |
| No-build mode | Yes | Yes |
| Nginx-only mode | Yes | Yes |
| Build time (all) | ~15-20 min | ~5-6 min |
| Build time (single) | ~15-20 min | ~1-2 min |

## Implementation Details

- **Inline Dockerfiles**: Uses `docker build -f - .` with heredoc to avoid file creation
- **Parallel execution**: Builds run in background with `&` and `wait`
- **Timing output**: Each step reports elapsed time for performance tracking
- **Smart skipping**: Only builds/transfers selected services
- **Backward compatible**: All nginx, certbot, and remote deployment steps unchanged
