#!/bin/bash
# Enable nginx sites by creating symlinks from sites-available to sites-enabled.
# Run as root or with sudo on the remote server.

set -e

NGINX_DIR="${NGINX_DIR:-/etc/nginx}"
SITES_AVAILABLE="$NGINX_DIR/sites-available"
SITES_ENABLED="$NGINX_DIR/sites-enabled"

echo "Enabling InsureTech nginx sites..."

# All sites to enable
SITES=(
    "default.conf"
    "insuretech-api.conf"
    "b2b_portal.conf"
    "labaidinsuretech.com.conf"
    "system_portal.conf"
    "coming_soon.conf"
)

for SITE in "${SITES[@]}"; do
    SRC="$SITES_AVAILABLE/$SITE"
    DEST="$SITES_ENABLED/$SITE"

    if [ ! -f "$SRC" ]; then
        echo "  WARNING: $SRC not found, skipping."
        continue
    fi

    # Remove stale symlink if exists
    [ -L "$DEST" ] && rm "$DEST"
    # Remove regular file if somehow exists
    [ -f "$DEST" ] && sudo rm "$DEST"

    ln -s "$SRC" "$DEST"
    echo "  ✓ Enabled: $SITE"
done

# Remove old/stale configs that no longer exist
for LINK in "$SITES_ENABLED"/*.conf; do
    LINK_TARGET=$(readlink "$LINK" 2>/dev/null || true)
    if [ -n "$LINK_TARGET" ] && [ ! -f "$LINK_TARGET" ]; then
        echo "  Removing stale symlink: $LINK"
        rm "$LINK"
    fi
done

# Clean up old domain names
rm -f "$SITES_ENABLED"/trendyco*.conf 2>/dev/null || true

echo "Testing nginx configuration..."
nginx -t

echo "Reloading nginx..."
systemctl reload nginx 2>/dev/null || nginx -s reload

echo "Done! Active sites:"
ls -la "$SITES_ENABLED/"
