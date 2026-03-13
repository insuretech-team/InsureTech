#!/bin/bash
# =============================================================
# InsureTech -- Quick Docker Deployment Script
# =============================================================
# Builds images locally, transfers to remote, loads + starts
# the full stack, syncs nginx, and provisions SSL (first run).
#
# Usage:
#   bash scripts/quick_deploy.sh              # full deploy
#   bash scripts/quick_deploy.sh --nginx-only # nginx + certs only
#   bash scripts/quick_deploy.sh --no-build   # skip build (cached)
# =============================================================

set -euo pipefail

# -- Config ----------------------------------------------------
REMOTE_HOST="insureadmin@146.190.97.242"
REMOTE_DIR="/home/insureadmin/insuretech"
NGINX_CONF_DIR="backend/infra/nginx"
CERTBOT_EMAIL="admin@labaidinsuretech.com"

# -- Flags -----------------------------------------------------
NGINX_ONLY=false
NO_BUILD=false
for arg in "$@"; do
    case $arg in
        --nginx-only) NGINX_ONLY=true ;;
        --no-build)   NO_BUILD=true ;;
    esac
done

# -- Docker binary (WSL: prefer docker, fall back to docker.exe) --
if command -v docker &>/dev/null; then
    DOCKER="docker"
elif command -v docker.exe &>/dev/null; then
    DOCKER="docker.exe"
else
    echo "ERROR: docker not found. Enable Docker Desktop WSL integration."
    exit 1
fi

echo "=== InsureTech Quick Docker Deployment ==="
echo "  Target : $REMOTE_HOST"
echo "  Mode   : $([ "$NGINX_ONLY" = true ] && echo 'nginx+certs only' \
    || ([ "$NO_BUILD" = true ] && echo 'no-build (cached images)' \
    || echo 'full deploy'))"
echo ""

# -- SSH Multiplexing (ControlMaster) --------------------------
# Opens ONE master SSH connection upfront. All subsequent ssh/scp
# calls reuse this socket -- zero additional password prompts.
SSH_SOCKET="/tmp/insuretech_ssh_$$"
# ControlPersist=7200 = 2 hours -- long enough to survive full image build+transfer.
# StrictHostKeyChecking=no avoids interactive fingerprint prompts.
SSH_OPTS="-o ControlMaster=auto -o ControlPath=$SSH_SOCKET -o ControlPersist=7200 -o StrictHostKeyChecking=no"

# shellcheck disable=SC2064
trap "ssh -O exit -o ControlPath=$SSH_SOCKET $REMOTE_HOST 2>/dev/null || true; rm -f $SSH_SOCKET" EXIT

echo ">> Establishing SSH master connection to $REMOTE_HOST ..."
echo "   (You will be prompted for SSH password ONCE)"
# -N = no remote command, just open the master connection
ssh $SSH_OPTS -N -f "$REMOTE_HOST"
echo "   SSH master connection established."
echo ""

# Wrappers that reuse the master socket -- no more prompts
ssh_run() { ssh $SSH_OPTS "$REMOTE_HOST" "$@"; }
scp_put() { scp -o ControlMaster=auto -o ControlPath="$SSH_SOCKET" -o StrictHostKeyChecking=no "$@"; }

# Running as root -- no sudo password needed.
SUDO_PASS=""

# -- Load .env.prod --------------------------------------------
# Use set -a + source (NOT xargs) -- .env.prod has multi-line PEM
# keys that xargs breaks into invalid identifiers.
if [ -f ".env.prod" ]; then
    set -a
    source .env.prod
    set +a
fi
B2B_BUILD_URL="${B2B_PORTAL_API_BASE_URL:-http://146.190.97.242}"
echo ">> B2B Portal API URL (baked into Next.js): $B2B_BUILD_URL"
echo ""

# =============================================================
# STEP 0 -- Push .env.prod to remote FIRST
# CRITICAL: db-migrate uses env_file: .env on remote.
# docker compose up runs AFTER this step so db-migrate gets
# the correct Neon DB credentials, not dev localhost creds.
# =============================================================
echo "[0/4] Syncing .env.prod to remote (MUST happen before stack start)..."
scp_put ".env.prod" "$REMOTE_HOST:$REMOTE_DIR/.env.prod"
scp_put ".env.prod" "$REMOTE_HOST:$REMOTE_DIR/.env"
echo "   Done."
echo ""

# =============================================================
# STEP 1 -- Build images locally & transfer to remote
# =============================================================
if [ "$NGINX_ONLY" = false ]; then
    echo "[1/4] Docker images..."

    if [ "$NO_BUILD" = false ]; then
        echo "  Building images locally..."
        # We use `docker build` directly (not `docker compose build`) because:
        # 1. `docker compose build` in BuildKit bake mode generates content-hash
        #    image names ignoring our explicit image: tags in compose.
        # 2. Building directly with -t gives us the exact insuretech-* tag we need.
        # 3. --no-cache busts any stale BuildKit layers.

        # Application microservices
        for SPEC in \
            "authn:backend/infra/docker/authn/Dockerfile:insuretech-authn" \
            "authz:backend/infra/docker/authz/Dockerfile:insuretech-authz" \
            "b2b:backend/infra/docker/b2b/Dockerfile:insuretech-b2b" \
            "gateway:backend/infra/docker/gateway/Dockerfile:insuretech-gateway" \
            "storage:backend/infra/docker/storage/Dockerfile:insuretech-storage" \
            "dbops:backend/infra/docker/dbops/Dockerfile:insuretech-dbops"; do
            SVC=$(echo "$SPEC" | cut -d: -f1)
            DFILE=$(echo "$SPEC" | cut -d: -f2)
            TAG=$(echo "$SPEC" | cut -d: -f3)
            echo ""
            echo "  ======================================================"
            echo "    Building $SVC -> $TAG ..."
            echo "  ======================================================"
            $DOCKER build --no-cache --progress=plain -f "$DFILE" -t "$TAG" . 2>&1
            echo "  BUILD OK: $TAG"
        done

        # b2b_portal (Next.js) -- needs API URL build arg
        echo ""
        echo "  ======================================================"
        echo "    Building b2b_portal -> insuretech-b2b_portal ..."
        echo "  ======================================================"
        $DOCKER build --no-cache --progress=plain \
            -f backend/infra/docker/b2b_portal/Dockerfile \
            --build-arg NEXT_PUBLIC_INSURETECH_API_BASE_URL="$B2B_BUILD_URL" \
            -t insuretech-b2b_portal . 2>&1
        echo "  BUILD OK: insuretech-b2b_portal"

        echo "  Pulling infra images..."
        $DOCKER pull redis:7-alpine      || true
        $DOCKER pull apache/kafka:latest || true
    fi

    echo "  Saving images to tarball..."
    mkdir -p build
    IMAGE_TAR="build/insuretech_images.tar"

    # docker-compose-prod.yml is the production compose file.
    # We use it locally to enumerate images, then push it to remote AS docker-compose.yml.
    COMPOSE_IMAGES=$($DOCKER compose -f docker-compose-prod.yml --profile full config --images)
    $DOCKER save -o "$IMAGE_TAR" $COMPOSE_IMAGES
    echo "  Tarball: $(du -sh "$IMAGE_TAR" | cut -f1)"

    echo "  Transferring docker-compose-prod.yml to remote as docker-compose.yml..."
    scp_put "docker-compose-prod.yml" "$REMOTE_HOST:$REMOTE_DIR/docker-compose.yml"

    echo "  Transferring image tarball to remote (may take a few minutes)..."
    scp_put "$IMAGE_TAR" "$REMOTE_HOST:$REMOTE_DIR/insuretech_images.tar"

    echo "  Transferring b2b_portal/public/ assets to remote..."
    # Pack public/ into a tar locally, ship one file, extract on remote.
    # This avoids scp -r quirks with trailing /. on Windows SSH clients
    # and handles subdirectories (logos/, stats-cards/, icons/ etc.) correctly.
    tar -czf build/b2b_portal_public.tar.gz -C b2b_portal/public .
    ssh_run "rm -rf $REMOTE_DIR/b2b_portal_public && mkdir -p $REMOTE_DIR/b2b_portal_public"
    scp_put "build/b2b_portal_public.tar.gz" "$REMOTE_HOST:$REMOTE_DIR/b2b_portal_public.tar.gz"
    ssh_run "tar -xzf $REMOTE_DIR/b2b_portal_public.tar.gz -C $REMOTE_DIR/b2b_portal_public && rm -f $REMOTE_DIR/b2b_portal_public.tar.gz"
    rm -f build/b2b_portal_public.tar.gz
    echo "  b2b_portal/public/ transferred ($(find b2b_portal/public -type f | wc -l) files)."
    echo ""
fi

# =============================================================
# STEP 2 -- Sync nginx configs to remote
# =============================================================
echo "[2/4] Syncing nginx configuration..."

# Create all required directories on remote
ssh_run "sudo mkdir -p \
    /etc/nginx/sites-available /etc/nginx/sites-enabled \
    /etc/nginx/snippets /etc/nginx/conf.d /etc/nginx/upstreams \
    /etc/nginx/cache /etc/nginx/maps /etc/nginx/error-pages \
    /var/www/html /var/www/certbot /var/www/static \
    /var/www/labaidinsuretech.com /var/www/b2b.labaidinsuretech.com \
    /var/www/system.labaidinsuretech.com /var/www/coming_soon \
    /var/cache/nginx/static /var/cache/nginx/api /var/cache/nginx/microcache"

# Copy configs via /tmp (scp to user-writable dir, then sudo mv to /etc/nginx)
# We copy sites-available only -- NOT sites-enabled.
# sites-enabled is managed exclusively by the symlink step (section 3f)
# to avoid duplicate server_name warnings.
ssh_run "rm -rf /tmp/nginx_upload && mkdir -p /tmp/nginx_upload"
scp_put -r "$NGINX_CONF_DIR/." "$REMOTE_HOST:/tmp/nginx_upload/"
ssh_run "sudo bash -c '
    # Clear sites-enabled completely -- will be repopulated by symlinks below
    rm -f /etc/nginx/sites-enabled/*.conf
    # Remove stale configs that are not part of our deployment
    rm -f /etc/nginx/sites-available/customer_portal.conf
    rm -f /etc/nginx/sites-enabled/customer_portal.conf
    # Copy everything except sites-enabled
    find /tmp/nginx_upload -mindepth 1 -maxdepth 1 -not -name sites-enabled \
        -exec cp -r {} /etc/nginx/ \;
    rm -rf /tmp/nginx_upload
    echo Done copying nginx configs
'"
echo "   Nginx configs synced."
echo ""

# =============================================================
# STEP 3 -- Run everything on remote in ONE SSH session
# =============================================================
echo "[3/4] Running remote deployment..."

# We embed SUDO_PASS directly into the script string so remote
# sudo -S can read it from stdin without a TTY.
# NOTE: heredoc uses unquoted delimiter so local vars (${SUDO_PASS},
# ${NGINX_ONLY}, ${CERTBOT_EMAIL}, ${REMOTE_DIR}) expand NOW (local side),
# while \$-prefixed vars expand LATER (remote side).
ssh_run "bash -s" << REMOTE_EOF
set -euo pipefail
SUDO_PASS='${SUDO_PASS}'
NGINX_ONLY='${NGINX_ONLY}'
CERTBOT_EMAIL='${CERTBOT_EMAIL}'
REMOTE_DIR='${REMOTE_DIR}'

# insureadmin has NOPASSWD:ALL -- use sudo for privileged ops
s() {
    sudo "\$@"
}

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

    echo "  Loading Docker images from tarball..."
    docker load -i insuretech_images.tar
    rm -f insuretech_images.tar

    echo "  Starting Docker stack (--profile full)..."
    docker compose --profile full down --remove-orphans 2>/dev/null || true
    docker compose --profile full up -d --no-build --remove-orphans

    echo "  Waiting 30s for containers to initialise..."
    sleep 30
    docker compose --profile full ps

    # -- 3b-ii. Sync b2b_portal public/ assets into running container ----
    # The public/ folder (logos, favicons, static assets) is baked into the
    # Docker image at build time. However, assets updated locally after the
    # last full image build will be stale in the container unless we patch
    # them here. This step is fast and idempotent.
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
else
    # --nginx-only: no image transfer, but still apply any .env changes
    # by force-recreating containers that depend on env vars (gateway, b2b_portal).
    cd "\$REMOTE_DIR"
    echo "  Applying .env changes to gateway and b2b_portal (force-recreate)..."
    docker compose --profile full up -d --no-build --force-recreate gateway b2b_portal 2>/dev/null || true
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

# Ensure nginx.conf includes our upstreams dir (idempotent)
if ! grep -q "include /etc/nginx/upstreams/\*.conf;" /etc/nginx/nginx.conf 2>/dev/null; then
    s sed -i '/http {/a\\    include /etc/nginx/upstreams/*.conf;' /etc/nginx/nginx.conf
fi

# -- Remove certbot-added duplicate redirect server blocks ----
# Certbot adds a server block like:
#   server {
#       if (\$host = domain.com) { return 301 https://...; } # managed by Certbot
#       listen 80;
#       server_name domain.com;
#       return 404; # managed by Certbot
#   }
# This creates ERR_TOO_MANY_REDIRECTS when combined with our own HTTP block.
# We strip these certbot-added blocks after each deploy.
# NOTE: We write the Python script to /tmp first (not a nested heredoc) so that
# sudo -S can read the password from stdin without conflicting with the script body.
echo "  Removing certbot duplicate redirect blocks..."
cat > /tmp/strip_certbot.py << 'PYEOF'
import sys, re
path = sys.argv[1]
with open(path) as f:
    content = f.read()
cleaned = re.sub(
    r'\n*server\s*\{[^{}]*?managed by Certbot[^{}]*?return 404[^{}]*?\}\s*',
    '\n',
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

# -- 3e. Per-domain SSL enable/disable based on cert existence --
# For each conf file: if the cert for its primary domain exists,
# UNCOMMENT ssl directives. If not, comment them out.
# This is idempotent -- safe to run on every deploy.
echo "  Syncing SSL directives per domain cert status..."
declare -A DOMAIN_CERT_MAP
DOMAIN_CERT_MAP=(
    [b2b_portal.conf]="b2b.labaidinsuretech.com"
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
            -e 's|^#\(\s*listen 443 ssl\)|\1|g' \
            -e 's|^#\(\s*listen \[::\]:443 ssl\)|\1|g' \
            -e 's|^#\(\s*ssl_certificate \)|\1|g' \
            -e 's|^#\(\s*ssl_certificate_key \)|\1|g' \
            -e 's|^#\(\s*ssl_trusted_certificate \)|\1|g' \
            "\$CONF_PATH" 2>/dev/null || true
    else
        echo "    SSL DISABLE: \$CONF_FILE (no cert for \$CERT_DOMAIN yet)"
        s sed -i \
            -e 's|^\(\s*\)listen 443 ssl|#\1listen 443 ssl|g' \
            -e 's|^\(\s*\)listen \[::\]:443 ssl|#\1listen [::]:443 ssl|g' \
            -e 's|^\(\s*\)ssl_certificate |#\1ssl_certificate |g' \
            -e 's|^#\?\(\s*\)ssl_certificate_key |#\1ssl_certificate_key |g' \
            -e 's|^\(\s*\)ssl_trusted_certificate |#\1ssl_trusted_certificate |g' \
            "\$CONF_PATH" 2>/dev/null || true
    fi
done

# -- 3f. Enable sites -----------------------------------------
# Enable ALL sites from sites-available dynamically.
# sites-enabled was cleared before the config copy so there are
# no stale or duplicate configs here.
echo "  Enabling nginx sites..."
for SRC in /etc/nginx/sites-available/*.conf; do
    SITE=\$(basename "\$SRC")
    LINK="/etc/nginx/sites-enabled/\$SITE"
    s ln -sf "\$SRC" "\$LINK"
    echo "    ok: \$SITE"
done

# -- 3g. Firewall ---------------------------------------------
s ufw allow 'Nginx Full' 2>/dev/null || true

# -- 3h. Test + restart nginx ---------------------------------
echo "  Testing nginx config..."
s nginx -t
s systemctl restart nginx
echo "  nginx restarted ok."

# -- 3i. Certbot SSL (first run only) -------------------------
echo "  Checking SSL certificates..."
ALL_CERTIFIED=true
for CHECK_D in labaidinsuretech.com api.labaidinsuretech.com b2b.labaidinsuretech.com \
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

# =============================================================
# STEP 4 -- Done
# =============================================================
echo ""
echo "[4/4] Deployment complete!"
echo ""
echo "  Direct (debug) : http://146.190.97.242:8080/healthz  (gateway)"
echo "                   http://146.190.97.242:3000         (b2b portal)"
echo "  Via Nginx HTTPS: https://api.labaidinsuretech.com   -> gateway:8080"
echo "                   https://b2b.labaidinsuretech.com   -> portal:3000"
echo "                   https://api.labaidinsuretech.com/nginx-health (nginx liveness)"
echo ""
echo "  NOTE: AUTHN_GRPC_ADDR, AUTHZ_GRPC_ADDR, B2B_GRPC_ADDR, STORAGE_GRPC_ADDR"
echo "        are set in .env.prod so gateway resolves Docker service hostnames correctly."
echo "  Once DNS points to 146.190.97.242, re-run --nginx-only to get SSL certs."
