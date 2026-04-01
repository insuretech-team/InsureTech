#!/bin/bash
# =============================================================
# InsureTech -- Quicker Docker Deployment Script
# =============================================================
# FASTER version of quick_deploy.sh that:
# 1. Builds Go binaries NATIVELY in WSL (not inside Docker)
# 2. Uses slim alpine Dockerfiles with pre-built binaries
# 3. Accepts --services flag to build only selected services
#
# Usage:
#   bash scripts/quickerdeploy.sh                                    # full deploy (all services)
#   bash scripts/quickerdeploy.sh --services=authn,gateway,b2b_portal # selective deploy
#   bash scripts/quickerdeploy.sh --skip-build                       # skip build, use existing images, still save/transfer/load
#   bash scripts/quickerdeploy.sh --restart-only (or --no-build)     # restart only, no save/transfer/load
#   bash scripts/quickerdeploy.sh --nginx-only                       # nginx + certs only
# =============================================================

set -euo pipefail

# -- Config ----------------------------------------------------
REMOTE_HOST="insureadmin@146.190.97.242"
REMOTE_DIR="/home/insureadmin/insuretech"
NGINX_CONF_DIR="backend/infra/nginx"
CERTBOT_EMAIL="admin@labaidinsuretech.com"

# -- Flags and Service Selection --------------------------------
NGINX_ONLY=false
NO_BUILD=false
SKIP_BUILD=false
SERVICES_ARG=""
SERVICES=(authn authz b2b kyc partner gateway storage dbops workflow payment fraud notification media insurance orders docgen polisync docrender b2b_portal insurer_portal)

for arg in "$@"; do
    case $arg in
        --nginx-only)     NGINX_ONLY=true ;;
        --skip-build)     SKIP_BUILD=true ;;
        --no-build|--restart-only)  NO_BUILD=true ;;
        --services=*)     SERVICES_ARG="${arg#--services=}" ;;
    esac
done

# Parse --services comma-separated list
if [ -n "$SERVICES_ARG" ]; then
    IFS=',' read -ra SERVICES <<< "$SERVICES_ARG"
fi

# -- Helper function to check if service is in SERVICES array --
contains_service() {
    local service="$1"
    for s in "${SERVICES[@]}"; do
        [ "$s" = "$service" ] && return 0
    done
    return 1
}

# -- Docker binary (WSL: prefer docker, fall back to docker.exe) --
if command -v docker &>/dev/null; then
    DOCKER="docker"
elif command -v docker.exe &>/dev/null; then
    DOCKER="docker.exe"
else
    echo "ERROR: docker not found. Enable Docker Desktop WSL integration."
    exit 1
fi

echo "=== InsureTech Quicker Docker Deployment ==="
echo "  Target   : $REMOTE_HOST"
echo "  Services : ${SERVICES[*]}"
echo "  Mode     : $([ "$NGINX_ONLY" = true ] && echo 'nginx+certs only' \
    || ([ "$NO_BUILD" = true ] && echo 'restart-only (no save/transfer/load)' \
    || ([ "$SKIP_BUILD" = true ] && echo 'skip-build (use existing images, still save/transfer/load)' \
    || echo 'full deploy (native build + slim Docker)')))"
echo ""

# -- SSH Multiplexing (ControlMaster) --------------------------
SSH_SOCKET="/tmp/insuretech_ssh_$$"
SSH_OPTS="-o ControlMaster=auto -o ControlPath=$SSH_SOCKET -o ControlPersist=7200 -o StrictHostKeyChecking=no"

# shellcheck disable=SC2064
trap "ssh -O exit -o ControlPath=$SSH_SOCKET $REMOTE_HOST 2>/dev/null || true; rm -f $SSH_SOCKET; rm -f .insuretech_dockerfile.tmp" EXIT

echo ">> Establishing SSH master connection to $REMOTE_HOST ..."
echo "   (You will be prompted for SSH password ONCE)"
ssh $SSH_OPTS -N -f "$REMOTE_HOST"
echo "   SSH master connection established."
echo ""

# Wrappers that reuse the master socket
ssh_run() { ssh $SSH_OPTS "$REMOTE_HOST" "$@"; }
scp_put() { scp -o ControlMaster=auto -o ControlPath="$SSH_SOCKET" -o StrictHostKeyChecking=no "$@"; }

# -- Load .env.prod (save flags first so .env.prod cannot override them) -------
if [ -f ".env.prod" ]; then
    _NGINX_ONLY="$NGINX_ONLY"
    _NO_BUILD="$NO_BUILD"
    _SKIP_BUILD="$SKIP_BUILD"
    set -a
    source .env.prod
    set +a
    # Restore deploy flags — .env.prod must not override CLI flags
    NGINX_ONLY="$_NGINX_ONLY"
    NO_BUILD="$_NO_BUILD"
    SKIP_BUILD="$_SKIP_BUILD"
fi
B2B_BUILD_URL="${B2B_PORTAL_API_BASE_URL:-http://146.190.97.242}"
INSURER_BUILD_URL="${INSURER_PORTAL_API_BASE_URL:-http://146.190.97.242}"
echo ">> B2B Portal API URL (baked into Next.js): $B2B_BUILD_URL"
echo ">> Insurer Portal API URL (baked into Next.js): $INSURER_BUILD_URL"
echo ""

# =============================================================
# STEP 0 -- Push .env.prod to remote FIRST
# =============================================================
echo "[0/5] Syncing .env.prod, nginx configs, and secrets to remote..."
scp_put ".env.prod" "$REMOTE_HOST:$REMOTE_DIR/.env.prod"
scp_put ".env.prod" "$REMOTE_HOST:$REMOTE_DIR/.env"

# Always sync secrets — needed for authz/polisync volume mounts regardless of build mode
echo "  Syncing backend/inscore/secrets (JWT keys etc.)..."
ssh_run "mkdir -p $REMOTE_DIR/backend/inscore/secrets"
scp_put -r "backend/inscore/secrets/." "$REMOTE_HOST:$REMOTE_DIR/backend/inscore/secrets/"

# Render nginx conf.d templates locally from .env.prod before syncing
# DEPLOY_DIR is the project root (where quickerdeploy.sh must be run from)
DEPLOY_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
NGINX_INFRA_LOCAL="$DEPLOY_DIR/backend/infra/nginx"
ENV_PROD_FILE="$DEPLOY_DIR/.env.prod"

echo "  Rendering nginx conf.d templates from $ENV_PROD_FILE..."
bash "$NGINX_INFRA_LOCAL/scripts/render-nginx-conf.sh" \
    --env "$ENV_PROD_FILE" \
    --out "$NGINX_INFRA_LOCAL/dist/conf.d" 2>&1 | sed 's/^/  /'

# Sync rendered conf.d (env-substituted) and sites-available to remote
ssh_run "mkdir -p $REMOTE_DIR/backend/infra/nginx/dist/conf.d $REMOTE_DIR/backend/infra/nginx/sites-available"
scp_put -r "${NGINX_INFRA_LOCAL}/dist/conf.d/."        "$REMOTE_HOST:$REMOTE_DIR/backend/infra/nginx/dist/conf.d/"
scp_put -r "${NGINX_INFRA_LOCAL}/sites-available/."    "$REMOTE_HOST:$REMOTE_DIR/backend/infra/nginx/sites-available/"
echo "   Done."
echo ""

# =============================================================
# STEP 1-3 -- Build locally (skip if --skip-build, --no-build, or --nginx-only)
# =============================================================
if [ "$NGINX_ONLY" = false ] && [ "$NO_BUILD" = false ]; then

if [ "$SKIP_BUILD" = false ]; then
    # ========================================================
    # STEP 1: Nuke stale binaries for selected services ONLY
    # ========================================================
    echo "[1/5] Removing stale binaries for selected services..."
    START_TIME=$(date +%s)
    for service in authn authz b2b kyc partner gateway storage dbops workflow payment fraud notification media insurance orders docgen; do
        if contains_service "$service"; then
            rm -rf "bin/$service"
            echo "  removed: bin/$service"
        fi
    done
    echo "   Done in $(($(date +%s) - START_TIME))s"
    echo ""

    # ========================================================
    # STEP 2: Build Go binaries sequentially (reliable output)
    # ========================================================
    echo "[2/5] Building Go binaries (CGO_ENABLED=0, linux/amd64)..."
    BUILD_START=$(date +%s)

    declare -A SERVICE_SOURCE_PATHS=(
        [authn]="./backend/inscore/microservices/authn/cmd/server/main.go"
        [authz]="./backend/inscore/microservices/authz/cmd/server/main.go"
        [b2b]="./backend/inscore/microservices/b2b/cmd/server/main.go"
        [kyc]="./backend/inscore/microservices/kyc/cmd/server/main.go"
        [partner]="./backend/inscore/microservices/partner/cmd/server/main.go"
        [gateway]="./backend/inscore/cmd/gateway/main.go"
        [storage]="./backend/inscore/cmd/storage/main.go"
        [dbops]="./backend/inscore/cmd/dbops/main.go"
        [workflow]="./backend/inscore/microservices/workflow/cmd/server/main.go"
        [payment]="./backend/inscore/cmd/payment/main.go"
        [fraud]="./backend/inscore/microservices/fraud/cmd/server/main.go"
        [notification]="./backend/inscore/microservices/notification/cmd/server/main.go"
        [media]="./backend/inscore/microservices/media/cmd/server/main.go"
        [insurance]="./backend/inscore/microservices/insurance/cmd/server/main.go"
        [orders]="./backend/inscore/microservices/orders/cmd/server/main.go"
        [docgen]="./backend/inscore/microservices/docgen/cmd/server/main.go"
    )

    for service in authn authz b2b kyc partner gateway storage dbops workflow payment fraud notification media insurance orders docgen; do
        if ! contains_service "$service"; then
            continue
        fi
        src="${SERVICE_SOURCE_PATHS[$service]}"
        outdir="bin/$service"
        outbin="$outdir/server"
        [ "$service" = "dbops" ] && outbin="$outdir/dbops"
        mkdir -p "$outdir"
        svc_start=$(date +%s)
        echo "  [$service] Building $src ..."
        if CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
            -ldflags="-w -s" -o "$outbin" "$src"; then
            echo "  [$service] ✓ $(du -sh "$outbin" | cut -f1) — $(($(date +%s) - svc_start))s"
        else
            echo "  [$service] ✗ BUILD FAILED — aborting"
            exit 1
        fi
    done

    echo "   Total Go build time: $(($(date +%s) - BUILD_START))s"
    echo ""

    # ========================================================
    # STEP 3: Build Docker images using temporary Dockerfile files
    # (avoids heredoc quoting issues on all platforms)
    # ========================================================
    echo "[3/5] Building Docker images (slim alpine + pre-built binaries)..."
    DOCKER_BUILD_START=$(date +%s)
    TMPDF=".insuretech_dockerfile.tmp"

    for service in authn authz b2b kyc partner gateway storage workflow payment fraud notification media insurance orders; do
        if ! contains_service "$service"; then
            continue
        fi
        svc_start=$(date +%s)
        echo ""
        echo "  Building insuretech-$service ..."
        cat > "$TMPDF" <<EODF
FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY bin/${service}/server .
COPY go.mod .
COPY ops ./ops
COPY backend/inscore/configs ./backend/inscore/configs
COPY backend/inscore/configs ./inscore/configs
COPY backend/inscore/templates ./backend/inscore/templates
COPY backend/inscore/secrets ./backend/inscore/secrets
CMD ["./server"]
EODF
        $DOCKER build --progress=plain -t "insuretech-$service" -f "$TMPDF" .
        echo "  ✓ insuretech-$service — $(($(date +%s) - svc_start))s"
    done

    # docgen — thin image with extra template/generated assets
    if contains_service "docgen"; then
        svc_start=$(date +%s)
        echo ""
        echo "  Building insuretech-docgen ..."
        $DOCKER build --progress=plain -t "insuretech-docgen" \
            -f backend/infra/docker/docgen/Dockerfile .
        echo "  ✓ insuretech-docgen — $(($(date +%s) - svc_start))s"
    fi

    # docrender — Python sidecar (must be built via Docker, not WSL)
    if contains_service "docrender"; then
        svc_start=$(date +%s)
        echo ""
        echo "  Building insuretech-docrender (Python sidecar) ..."
        $DOCKER build --progress=plain \
            -f backend/inscore/microservices/docgen/sidecar/Dockerfile \
            -t insuretech-docrender \
            backend/inscore/microservices/docgen/sidecar
        echo "  ✓ insuretech-docrender — $(($(date +%s) - svc_start))s"
    fi

    # polisync — C# .NET 8 multi-stage build (no WSL pre-compile possible)
    if contains_service "polisync"; then
        svc_start=$(date +%s)
        echo ""
        echo "  Building insuretech-polisync (C# .NET 8 multi-stage) ..."
        $DOCKER build --progress=plain \
            -f backend/polisync/Dockerfile \
            -t insuretech-polisync .
        echo "  ✓ insuretech-polisync — $(($(date +%s) - svc_start))s"
    fi

    # gotenberg — public image, pull on remote directly (do not bundle in tarball)
    # Handled in the remote deploy step below

    if contains_service "dbops"; then
        svc_start=$(date +%s)
        echo ""
        echo "  Building insuretech-dbops ..."
        cat > "$TMPDF" <<EODF
FROM alpine:latest
WORKDIR /app
COPY bin/dbops/dbops .
COPY go.mod .
COPY backend/inscore/db/migrations ./backend/inscore/db/migrations
COPY backend/inscore/db/seeds ./backend/inscore/db/seeds
COPY backend/inscore/configs ./backend/inscore/configs
COPY backend/inscore/secrets ./backend/inscore/secrets
CMD ["./dbops", "migrate", "--target=both"]
EODF
        $DOCKER build --progress=plain -t "insuretech-dbops" -f "$TMPDF" .
        echo "  ✓ insuretech-dbops — $(($(date +%s) - svc_start))s"
    fi

    rm -f "$TMPDF"

    if contains_service "b2b_portal"; then
        svc_start=$(date +%s)
        echo ""
        echo "  Building insuretech-b2b_portal (Node.js multi-stage) ..."
        FLVE_TOKEN="${NEXT_PUBLIC_FLVE_API_TOKEN:-${FLVE_API_TOKEN:-}}"
        $DOCKER build --progress=plain \
            -f backend/infra/docker/b2b_portal/Dockerfile \
            --build-arg NEXT_PUBLIC_INSURETECH_API_BASE_URL="$B2B_BUILD_URL" \
            --build-arg NEXT_PUBLIC_FLVE_API_TOKEN="$FLVE_TOKEN" \
            -t insuretech-b2b_portal .
        echo "  ✓ insuretech-b2b_portal — $(($(date +%s) - svc_start))s"
    fi

    if contains_service "insurer_portal"; then
        svc_start=$(date +%s)
        echo ""
        echo "  Building insuretech-insurer_portal (Node.js multi-stage) ..."
        $DOCKER build --progress=plain \
            -f backend/infra/docker/insurer_portal/Dockerfile \
            --build-arg NEXT_PUBLIC_INSURETECH_API_BASE_URL="$INSURER_BUILD_URL" \
            -t insuretech-insurer_portal .
        echo "  ✓ insuretech-insurer_portal — $(($(date +%s) - svc_start))s"
    fi

    # Pull infra images only on full deploy
    if contains_service "authn" && contains_service "gateway"; then
        echo ""
        echo "  Pulling infra images (redis, kafka)..."
        $DOCKER pull redis:7-alpine || true
        $DOCKER pull apache/kafka:latest || true
    fi

    echo "   Total Docker build time: $(($(date +%s) - DOCKER_BUILD_START))s"
    echo ""

fi  # end SKIP_BUILD=false (steps 1-3)

# ========================================================
# STEP 4: Save images to tarball (runs for full AND --skip-build)
# ========================================================
echo "[4/5] Saving Docker images to tarball..."
SAVE_START=$(date +%s)
mkdir -p build
IMAGE_TAR="build/insuretech_images.tar"

# Build list of images to save (only selected services from SERVICES array)
IMAGES_TO_SAVE=""
for service in authn authz b2b kyc partner gateway storage dbops workflow payment fraud notification media insurance orders docgen docrender polisync b2b_portal insurer_portal; do
    if contains_service "$service"; then
        # Verify image actually exists locally before adding to save list
        if $DOCKER image inspect "insuretech-$service" &>/dev/null; then
            IMAGES_TO_SAVE="$IMAGES_TO_SAVE insuretech-$service"
        else
            echo "  WARN: insuretech-$service image not found locally — skipping"
        fi
    fi
done

# gotenberg is a public image — pulled directly on remote, not bundled in tarball

# Include infra images only on full deploy (all services selected)
if contains_service "authn" && contains_service "gateway" && contains_service "b2b"; then
    IMAGES_TO_SAVE="$IMAGES_TO_SAVE redis:7-alpine apache/kafka:latest"
fi

if [ -z "$IMAGES_TO_SAVE" ]; then
    echo "ERROR: No images to save. Build may have failed."
    exit 1
fi

echo "  Images to save:$IMAGES_TO_SAVE"
$DOCKER save -o "$IMAGE_TAR" $IMAGES_TO_SAVE
echo "   Tarball size: $(du -sh "$IMAGE_TAR" | cut -f1)"
echo "   Saved in $(($(date +%s) - SAVE_START))s"
echo ""

fi  # end NO_BUILD=false (steps 1-4)

# =============================================================
# STEP 5: Transfer and deploy (skip if --restart-only or --nginx-only)
# =============================================================

if [ "$NGINX_ONLY" = false ] && [ "$NO_BUILD" = false ]; then
    TRANSFER_START=$(date +%s)
    
    if [ "$SKIP_BUILD" = false ]; then
        echo "[5/5] Transferring to remote..."
    else
        echo "[4/4] Transferring to remote (--skip-build)..."
    fi
    
    echo "  Transferring docker-compose-prod.yml ..."
    scp_put "docker-compose-prod.yml" "$REMOTE_HOST:$REMOTE_DIR/docker-compose.yml"

    echo "  Transferring image tarball (may take a few minutes)..."
    scp_put "build/insuretech_images.tar" "$REMOTE_HOST:$REMOTE_DIR/insuretech_images.tar"
    
    echo "  Transferring b2b_portal/public/ assets..."
    tar -czf build/b2b_portal_public.tar.gz -C b2b_portal/public . 2>/dev/null || true
    ssh_run "rm -rf $REMOTE_DIR/b2b_portal_public && mkdir -p $REMOTE_DIR/b2b_portal_public"
    scp_put "build/b2b_portal_public.tar.gz" "$REMOTE_HOST:$REMOTE_DIR/b2b_portal_public.tar.gz"
    ssh_run "tar -xzf $REMOTE_DIR/b2b_portal_public.tar.gz -C $REMOTE_DIR/b2b_portal_public && rm -f $REMOTE_DIR/b2b_portal_public.tar.gz"
    rm -f build/b2b_portal_public.tar.gz

    echo "  Transferring insurer-portal/public/ assets..."
    tar -czf build/insurer_portal_public.tar.gz -C insurer-portal/public . 2>/dev/null || true
    ssh_run "rm -rf $REMOTE_DIR/insurer_portal_public && mkdir -p $REMOTE_DIR/insurer_portal_public"
    scp_put "build/insurer_portal_public.tar.gz" "$REMOTE_HOST:$REMOTE_DIR/insurer_portal_public.tar.gz"
    ssh_run "tar -xzf $REMOTE_DIR/insurer_portal_public.tar.gz -C $REMOTE_DIR/insurer_portal_public && rm -f $REMOTE_DIR/insurer_portal_public.tar.gz"
    rm -f build/insurer_portal_public.tar.gz
    
    echo "   Transfer complete in $(($(date +%s) - TRANSFER_START))s"
    echo ""
fi

# =============================================================
# Remote deployment (same as quick_deploy.sh)
# =============================================================
REMOTE_START=$(date +%s)
if [ "$NGINX_ONLY" = true ]; then
    STEP="2/5"
elif [ "$NO_BUILD" = true ]; then
    STEP="2/2"
elif [ "$SKIP_BUILD" = true ]; then
    STEP="4/4"
else
    STEP="5/5"
fi
echo "[$STEP] Running remote deployment..."

ssh_run "bash -s" << REMOTE_EOF
set -euo pipefail
NGINX_ONLY='${NGINX_ONLY}'
CERTBOT_EMAIL='${CERTBOT_EMAIL}'
REMOTE_DIR='${REMOTE_DIR}'
NO_BUILD='${NO_BUILD}'
SKIP_BUILD='${SKIP_BUILD}'
DEPLOY_SERVICES='${SERVICES[*]}'

s() { sudo "\$@"; }

# -- 3a. Stop legacy bare-metal systemd services ---------------
echo "  Stopping legacy services..."
for svc in insuretech-gateway insuretech-authn insuretech-authz \
           insuretech-tenant insuretech-b2b; do
    s systemctl stop    "\$svc" 2>/dev/null || true
    s systemctl disable "\$svc" 2>/dev/null || true
done
s systemctl daemon-reload 2>/dev/null || true

# -- 3b. Load images + start Docker stack ---------------------
if [ "\$NGINX_ONLY" = false ]; then
    cd "\$REMOTE_DIR"

    if [ "\$NO_BUILD" = false ] && [ -f "\$REMOTE_DIR/insuretech_images.tar" ]; then
        echo "  Loading Docker images from tarball..."
        docker load -i insuretech_images.tar
        rm -f insuretech_images.tar
        echo "  Docker images loaded."
    elif [ "\$NO_BUILD" = true ]; then
        echo "  --restart-only: using already-loaded images on server."
    fi

    # Map deploy service names to docker compose service names
    declare -A SVC_MAP=(
        [authn]="authn"
        [authz]="authz"
        [b2b]="b2b"
        [kyc]="kyc"
        [partner]="partner"
        [gateway]="gateway"
        [storage]="storage"
        [dbops]="db-migrate"
        [workflow]="workflow"
        [payment]="payment"
        [fraud]="fraud"
        [notification]="notification"
        [media]="media"
        [insurance]="insurance"
        [orders]="orders"
        [docgen]="docgen"
        [docrender]="docrender"
        [polisync]="polisync"
        [b2b_portal]="b2b_portal"
        [insurer_portal]="insurer_portal"
    )

    # Build list of compose services to restart
    COMPOSE_SVCS=""
    for svc in \$DEPLOY_SERVICES; do
        mapped="\${SVC_MAP[\$svc]:-}"
        [ -n "\$mapped" ] && COMPOSE_SVCS="\$COMPOSE_SVCS \$mapped"
    done
    # gotenberg is a public image — pull it on remote directly before starting
    if echo "\$DEPLOY_SERVICES" | grep -qE '(docgen|docrender)'; then
        echo "  Pulling gotenberg/gotenberg:8 on remote..."
        docker pull gotenberg/gotenberg:8
        COMPOSE_SVCS="\$COMPOSE_SVCS gotenberg"
    fi
    COMPOSE_SVCS="\$(echo \$COMPOSE_SVCS | xargs)"  # trim

    if [ -n "\$COMPOSE_SVCS" ]; then
        echo "  Stopping services first (ensures fresh image is picked up)..."
        docker compose --profile full stop \$COMPOSE_SVCS 2>/dev/null || true
        docker compose --profile full rm -f \$COMPOSE_SVCS 2>/dev/null || true

        echo "  Starting services with fresh images: \$COMPOSE_SVCS"
        docker compose --profile full up -d --no-build --pull never --force-recreate \$COMPOSE_SVCS

        echo "  Waiting for containers to be running (max 60s)..."
        WAIT_SECS=0
        while [ \$WAIT_SECS -lt 60 ]; do
            sleep 5
            WAIT_SECS=\$((WAIT_SECS + 5))
            ALL_UP=true
            for svc in \$COMPOSE_SVCS; do
                # Use docker ps to check container state — simpler and more reliable
                CNAME="insuretech-\$svc"
                STATE=\$(docker ps --filter "name=^\${CNAME}$" --format "{{.Status}}" 2>/dev/null || echo "")
                if [ -z "\$STATE" ]; then
                    ALL_UP=false
                    echo "  ... [\$svc] not running yet (\${WAIT_SECS}s)"
                fi
            done
            \$ALL_UP && break
        done

        echo ""
        echo "  ════════════════════════════════════════════════════"
        echo "  Container status:"
        docker compose --profile full ps \$COMPOSE_SVCS
        echo ""
        echo "  ════════════════════════════════════════════════════"
        echo "  Live logs (last 60 lines per service):"
        echo "  ════════════════════════════════════════════════════"
        for svc in \$COMPOSE_SVCS; do
            echo ""
            echo "  ┌── [\$svc] ────────────────────────────────────"
            docker logs --tail=60 --timestamps "insuretech-\$svc" 2>&1 | sed 's/^/  │ /' || \
            docker compose --profile full logs --tail=60 \$svc 2>&1 | sed 's/^/  │ /' || true
            echo "  └────────────────────────────────────────────────"
        done
        echo ""
    else
        echo "  No compose services to restart."
    fi

    # -- 3b-ii. Sync b2b_portal public/ assets into running container ----
    if echo "\$DEPLOY_SERVICES" | grep -q "b2b_portal"; then
        echo "  Syncing b2b_portal public/ assets into container..."
        if docker ps --format '{{.Names}}' | grep -q '^insuretech-b2b-portal$'; then
            if [ -d "\$REMOTE_DIR/b2b_portal_public" ]; then
                docker cp "\$REMOTE_DIR/b2b_portal_public/." insuretech-b2b-portal:/app/public/
                rm -rf "\$REMOTE_DIR/b2b_portal_public"
                echo "  b2b_portal public/ assets synced."
            fi
        else
            echo "  WARN: insuretech-b2b-portal not running -- skipping public/ sync."
        fi
    fi

    if echo "\$DEPLOY_SERVICES" | grep -q "insurer_portal"; then
        echo "  Syncing insurer_portal public/ assets into container..."
        if docker ps --format '{{.Names}}' | grep -q '^insuretech-insurer-portal$'; then
            if [ -d "\$REMOTE_DIR/insurer_portal_public" ]; then
                docker cp "\$REMOTE_DIR/insurer_portal_public/." insuretech-insurer-portal:/app/public/
                rm -rf "\$REMOTE_DIR/insurer_portal_public"
                echo "  insurer_portal public/ assets synced."
            fi
        else
            echo "  WARN: insuretech-insurer-portal not running -- skipping public/ sync."
        fi
    fi
else
    # --nginx-only: no image transfer, but still apply any .env changes
    cd "\$REMOTE_DIR"
    echo "  Applying .env changes to gateway and frontend portals (force-recreate)..."
    docker compose --profile full up -d --no-build --force-recreate gateway b2b_portal insurer_portal 2>/dev/null || true
fi

# -- 3c. Ensure nginx + certbot installed ---------------------
echo "  Checking nginx/certbot..."
if ! command -v nginx &>/dev/null; then
    s apt-get update -qq
    s apt-get install -yq nginx
fi
if ! command -v certbot &>/dev/null; then
    s apt-get install -yq certbot python3-certbot-nginx
fi
s systemctl enable nginx

# -- 3d. Configure nginx --------------------------------------
echo "  Configuring nginx..."
s rm -f /etc/nginx/sites-enabled/default
s rm -f /etc/nginx/sites-available/default
s rm -f /etc/nginx/sites-enabled/trendyco* 2>/dev/null || true

if ! grep -q "include /etc/nginx/upstreams/\*.conf;" /etc/nginx/nginx.conf 2>/dev/null; then
    s sed -i '/http {/a\\    include /etc/nginx/upstreams/*.conf;' /etc/nginx/nginx.conf
fi

# -- 3d-i. Install pre-rendered nginx conf.d (rendered locally from .env.prod)
# conf.d files were already rendered (placeholders substituted) by render-nginx-conf.sh
# in step [0/5] on the local machine before SCP. No sed substitution needed here.
echo "  Installing pre-rendered nginx conf.d..."
NGINX_INFRA_DIR="\$REMOTE_DIR/backend/infra/nginx"
RENDERED_CONF_DIR="\$NGINX_INFRA_DIR/dist/conf.d"
if [ -d "\$RENDERED_CONF_DIR" ]; then
    s mkdir -p /etc/nginx/conf.d
    for CONF in "\$RENDERED_CONF_DIR/"*.conf; do
        [ -f "\$CONF" ] || continue
        s cp "\$CONF" "/etc/nginx/conf.d/\$(basename \$CONF)"
        echo "    installed: \$(basename \$CONF)"
    done
else
    echo "  ⚠  dist/conf.d not found — skipping conf.d update"
fi

# -- 3d-ii. Sync nginx sites-available from local project -------------------
echo "  Syncing nginx sites-available from project..."
if [ -d "\$NGINX_INFRA_DIR/sites-available" ]; then
    s mkdir -p /etc/nginx/sites-available
    for CONF in "\$NGINX_INFRA_DIR/sites-available/"*.conf; do
        [ -f "\$CONF" ] || continue
        s cp "\$CONF" "/etc/nginx/sites-available/\$(basename \$CONF)"
        echo "    synced: \$(basename \$CONF)"
    done
fi

# -- 3e. Remove certbot-added duplicate redirect server blocks ----
echo "  Removing certbot duplicate redirect blocks..."
cat > /tmp/strip_certbot.py << 'PYEOF'
import sys, re
path = sys.argv[1]
with open(path) as f:
    content = f.read()
cleaned = re.sub(
    r'\\n*server\\s*\\{[^{}]*?managed by Certbot[^{}]*?return 404[^{}]*?\\}\\s*',
    '\\n',
    content,
    flags=re.DOTALL
)
if cleaned != content:
    with open(path, 'w') as f:
        f.write(cleaned)
    print(f"  cleaned: {path}")
PYEOF
for CONF in /etc/nginx/sites-available/*.conf; do
    s python3 /tmp/strip_certbot.py "\$CONF"
done
rm -f /tmp/strip_certbot.py

# -- 3f. Per-domain SSL enable/disable based on cert existence --
echo "  Syncing SSL directives per domain cert status..."
declare -A DOMAIN_CERT_MAP
DOMAIN_CERT_MAP=(
    [b2b_portal.conf]="b2b.labaidinsuretech.com"
    [insurer_portal.conf]="insurer.labaidinsuretech.com"
    [insuretech-api.conf]="api.labaidinsuretech.com"
    [system_portal.conf]="system.labaidinsuretech.com"
    [coming_soon.conf]="agents.labaidinsuretech.com"
    [labaidinsuretech.com.conf]="labaidinsuretech.com"
)
for CONF_FILE in "\${!DOMAIN_CERT_MAP[@]}"; do
    CERT_DOMAIN="\${DOMAIN_CERT_MAP[\$CONF_FILE]}"
    CONF_PATH="/etc/nginx/sites-available/\$CONF_FILE"
    if s test -d "/etc/letsencrypt/live/\$CERT_DOMAIN" 2>/dev/null; then
        echo "    SSL ENABLE: \$CONF_FILE (cert exists for \$CERT_DOMAIN)"
        s sed -i \
            -e 's|^#\\(\\s*listen 443 ssl\\)|\\1|g' \
            -e 's|^#\\(\\s*listen \\[::\\]:443 ssl\\)|\\1|g' \
            -e 's|^#\\(\\s*ssl_certificate \\)|\\1|g' \
            -e 's|^#\\(\\s*ssl_certificate_key \\)|\\1|g' \
            -e 's|^#\\(\\s*ssl_trusted_certificate \\)|\\1|g' \
            "\$CONF_PATH" 2>/dev/null || true
    else
        echo "    SSL DISABLE: \$CONF_FILE (no cert for \$CERT_DOMAIN yet)"
        s sed -i \
            -e 's|^\\(\\s*\\)listen 443 ssl|#\\1listen 443 ssl|g' \
            -e 's|^\\(\\s*\\)listen \\[::\\]:443 ssl|#\\1listen [::]:443 ssl|g' \
            -e 's|^\\(\\s*\\)ssl_certificate |#\\1ssl_certificate |g' \
            -e 's|^#\\?\\(\\s*\\)ssl_certificate_key |#\\1ssl_certificate_key |g' \
            -e 's|^\\(\\s*\\)ssl_trusted_certificate |#\\1ssl_trusted_certificate |g' \
            "\$CONF_PATH" 2>/dev/null || true
    fi
done

# -- 3g. Enable sites -----------------------------------------
echo "  Enabling nginx sites..."
for SRC in /etc/nginx/sites-available/*.conf; do
    SITE=\$(basename "\$SRC")
    LINK="/etc/nginx/sites-enabled/\$SITE"
    s ln -sf "\$SRC" "\$LINK"
    echo "    ok: \$SITE"
done

# -- 3h. Firewall ---------------------------------------------
s ufw allow 'Nginx Full' 2>/dev/null || true

# -- 3i. Test + restart nginx ---------------------------------
echo "  Testing nginx config..."
s nginx -t
s systemctl restart nginx
echo "  nginx restarted ok."

# -- 3j. Certbot SSL (first run only) -------------------------
echo "  Checking SSL certificates..."
ALL_CERTIFIED=true
for CHECK_D in labaidinsuretech.com api.labaidinsuretech.com b2b.labaidinsuretech.com \
               insurer.labaidinsuretech.com \
               system.labaidinsuretech.com agents.labaidinsuretech.com \
               regulator.labaidinsuretech.com partners.labaidinsuretech.com \
               business.labaidinsuretech.com; do
    if ! s test -d "/etc/letsencrypt/live/\$CHECK_D" 2>/dev/null; then
        ALL_CERTIFIED=false
        echo "  Missing cert for \$CHECK_D"
    fi
done

if [ "\$ALL_CERTIFIED" = true ]; then
    echo "  All SSL certs exist -- skipping. Run: sudo certbot renew"
else
    echo "  Provisioning SSL certificates via Certbot..."
    CERT_GROUPS=(
        "labaidinsuretech.com www.labaidinsuretech.com"
        "api.labaidinsuretech.com"
        "b2b.labaidinsuretech.com"
        "insurer.labaidinsuretech.com"
        "system.labaidinsuretech.com"
        "agents.labaidinsuretech.com"
        "regulator.labaidinsuretech.com"
        "partners.labaidinsuretech.com"
        "business.labaidinsuretech.com"
    )
    for CERT_GROUP in "\${CERT_GROUPS[@]}"; do
        D_FLAGS=""
        for d in \$CERT_GROUP; do
            D_FLAGS="\$D_FLAGS -d \$d"
        done
        MAIN_D=\$(echo "\$CERT_GROUP" | awk '{print \$1}')
        echo "  Certbot: \$MAIN_D ..."
        if s certbot certificates 2>/dev/null | grep -q "\$MAIN_D"; then
            echo "    cert exists -- skipping (\$MAIN_D)"
        else
            s certbot --nginx --non-interactive --agree-tos \
                -m "\$CERTBOT_EMAIL" --no-redirect \
                \$D_FLAGS \
                || echo "    WARN: certbot failed for \$MAIN_D (DNS may not point here yet)"
        fi
    done
    s nginx -t && s systemctl reload nginx
fi

echo ""
echo "  nginx status: \$(s systemctl is-active nginx 2>/dev/null || echo unknown)"
echo "  Remote deployment complete."
REMOTE_EOF

echo "   Remote deployment completed in $(($(date +%s) - REMOTE_START))s"
echo ""

# =============================================================
# Done
# =============================================================
if [ "$NGINX_ONLY" = true ]; then
    FINAL_STEP="2/5"
elif [ "$NO_BUILD" = true ]; then
    FINAL_STEP="2/2"
elif [ "$SKIP_BUILD" = true ]; then
    FINAL_STEP="4/4"
else
    FINAL_STEP="5/5"
fi
echo "[$FINAL_STEP] Deployment complete!"
echo ""
echo "  Direct (debug) : http://146.190.97.242:8080/healthz  (gateway)"
echo "                   http://146.190.97.242:3000         (b2b portal)"
echo "                   http://127.0.0.1:3002             (insurer portal on server)"
echo "  Via Nginx HTTPS: https://api.labaidinsuretech.com   -> gateway:8080"
echo "                   https://b2b.labaidinsuretech.com   -> portal:3000"
echo "                   https://insurer.labaidinsuretech.com -> portal:3002"
echo "                   https://api.labaidinsuretech.com/nginx-health (nginx liveness)"
echo ""
echo "  NOTE: AUTHN_GRPC_ADDR, AUTHZ_GRPC_ADDR, B2B_GRPC_ADDR, STORAGE_GRPC_ADDR"
echo "        are set in .env.prod so gateway resolves Docker service hostnames correctly."
echo "  Once DNS points to 146.190.97.242, re-run --nginx-only to get SSL certs."
